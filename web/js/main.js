'use strict';

const DEFAULT_ALPHABET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'

document.ready(() => {

    let token = null
    const credentials = {
        login: localStorage.getItem('login'),
        password: localStorage.getItem('password'),
    }
    if(!credentials.login || !credentials.password){
        location.href = '/'
        return
    }

    const accessLevel = document.getElementById('access-level')
    const loginSpan = document.getElementById('credentials-login')
    if(loginSpan)
    {
        loginSpan.innerText = credentials.login
    }
    // Основные элементы интерфейса
    const start = document.getElementById('start')
    const stop = document.getElementById('stop')
    const remain = document.getElementById('remain')
    const app = document.getElementById('app')

    // Форма ввода паролей
    const result = document.getElementById('result')
    const currentPhrase = document.querySelector('label[for=result]')

    // Блок секретной записки
    const secretNoteCard = document.getElementById('secret-note-card')

    // Установка паролей
    const passwordsCard = document.getElementById('passwords-form-card')
    const passwordsForm = document.getElementById('passwords-form')
    let passwordsFormFields = []
    const passwordsFormBtn = document.getElementById('passwords-save')


    // Отобразить форму запуска
    const hideApp = () => {
        app.classList.add('d-none')

        start.classList.remove('d-none')
        stop.classList.add('d-none')
    }
    // Отобразить форму ввода
    const showApp = () => {
        app.classList.remove('d-none')

        start.classList.add('d-none')
        stop.classList.remove('d-none')
    }



    // События

    let session = null
    start.addEventListener('click', () => {
        processSession()
            .then((results) => {
                hideApp()
                console.log('READY', results)
                const data = {
                    auth: credentials,
                    flows: results,
                }
                fetch('/process', {
                    method: 'POST',
                    body: JSON.stringify(data),
                })
                    .then(res => res.json())
                    .then(res => {
                        console.log(res)
                        goCheckSecret()
                    })
            })
    })

    let onSessionStopped = () => {}
    function processSession(){
        if(session){
            session = null
            try{
                onSessionStopped()
            }catch (e){}
        }
        hideApp()

        return new Promise((resolve, reject) => {
            onSessionStopped = () => reject()

            fetch('/get-passwords', {
                method: 'POST',
                body: JSON.stringify(credentials),
            })
                .then(res => res.json())
                .then(res => {
                    if(res.success){
                        session = newSession({
                            phrases: res.data,
                            callback: resolve,
                        })
                        showApp()
                        session.next()
                    }
                })
                .catch(reject)
        })
    }

    stop.addEventListener('click', () => {
        session = null
        hideApp()
        onSessionStopped()
    })

    result.addEventListener('keydown', e => session && session.track(e))
    result.addEventListener('keyup', e => session && session.track(e))



    // Создание сессии
    function newSession({phrases, callback}){
        let row = []

        const sess = {
            current: '',
            result: [],
            source: phrases,
            stack: [],

            isReady: () => {
                const phraseIsReady = sess.current === result.value
                const allButtonsReleased = sess.stack.length === 0
                return phraseIsReady && allButtonsReleased
            },

            track: (keyboardEvent) => {
                // если кнопку зажали
                if(keyboardEvent.repeat){
                    keyboardEvent.preventDefault()
                    keyboardEvent.stopPropagation()
                    return
                }
                // если это спецклавиша
                if(keyboardEvent.key.length !== 1){
                    return
                }

                let prevent = false
                if(keyboardEvent.type === 'keyup'){
                    if(sess.stack.indexOf(keyboardEvent.key) === -1){
                        prevent = true
                    }else{
                        sess.stack = sess.stack.filter(x => x !== keyboardEvent.key)
                    }
                }else{
                    const candidate = result.value + (
                        keyboardEvent.shiftKey ? keyboardEvent.key.toUpperCase() : keyboardEvent.key.toLowerCase()
                    )
                    if(sess.current.indexOf(candidate) !== 0){
                        prevent = true
                    }else{
                        sess.stack.push(keyboardEvent.key)
                    }
                    // console.log(result.value, keyboardEvent.key)
                    // console.log(candidate, sess.current)
                }

                if(prevent){
                    keyboardEvent.preventDefault()
                    keyboardEvent.stopPropagation()
                    console.log('invalid character')
                    return
                }
                // console.log('stack:', ...sess.stack)

                row.push({
                    key: keyboardEvent.key,
                    time: Math.round(keyboardEvent.timeStamp),
                    up: (keyboardEvent.type === 'keyup'),
                })

                if(sess.isReady()){
                    result.value = ''
                    result.disabled = true
                    sess.result.push({
                        flow: row,
                        phrase: sess.current,
                    })
                    row = []
                    sess.next()
                }
            },

            next: () => {

                result.disabled = false
                result.value = ''
                remain.innerText = sess.source.length
                sess.current = sess.source.shift()
                if(!sess.current){
                    sess.ready()
                }else{
                    result.focus()
                }

                currentPhrase.innerText = sess.current ? sess.current : ''
            },

            ready: () => callback(sess.result),
        }

        console.log('session started', sess)
        return sess
    }


    // Конфигурация
    function config(){
        return new Promise(resolve => {
            fetch('/config')
                .then(res => res.json())
                .then(res => {
                    if(!res.success || !res.config){
                        return resolve(null)
                    }
                    return resolve(res.config)
                })
                .catch(err => {
                    console.error(err)
                    resolve(null)
                })
        })
    }


    // Установка паролей
    let passwordsFormRendered = false
    function renderPasswordsForm(){
        if(passwordsFormRendered){
            return
        }
        passwordsFormRendered = true

        config().then(({passwords_count}) => {
            if(!passwords_count){
                alert('Кол-во паролей неизвестно')
                return
            }
            passwordsFormFields = []
            for(let i = 0; i < passwords_count; i++){
                const pw = document.createElement('input')
                pw.type = 'text'
                pw.classList.add('form-control')
                pw.required = true
                pw.placeholder = 'Пароль №' + (i + 1)

                const inputGroup = document.createElement('div')
                inputGroup.classList.add('input-group')
                inputGroup.append(pw)

                const formGroup = document.createElement('div')
                formGroup.classList.add('form-group', 'mb-2')
                formGroup.append(inputGroup)

                passwordsForm.append(formGroup)

                passwordsFormFields.push(pw)
            }
            passwordsForm.addEventListener('submit', e => {
                e.preventDefault()
                passwordsFormBtn.disabled = true
                const data = {
                    auth: credentials,
                    passwords: passwordsFormFields.map(field => field.value)
                }
                fetch('/user/set-passwords', {
                    method: 'POST',
                    body: JSON.stringify(data),
                })
                    .then(res => res.json())
                    .then(res => {
                        if(!res.success){
                            passwordsFormBtn.disabled = false
                            return alert(res.error ?? 'Непредвиденная ошибка')
                        }
                        goCheckSecret()
                    })
                    .catch(err => {
                        passwordsFormBtn.disabled = false
                        alert(err)
                    })
            })
            passwordsFormBtn.disabled = false
        })
    }


    // Секретная записка
    const secret = document.getElementById('secret-note')
    const showSecretBtn = document.getElementById('secret-show')
    const saveSecretBtn = document.getElementById('secret-save')

    let isEditing = false
    let hasSecret = null
    let needPasswords = null
    let needSamples = null
    function checkSecret(){
        showSecretBtn.disabled = true
        return fetch('/user/has-secret', {
            method: 'POST',
            body: JSON.stringify(credentials),
        }).then(res => res.json())
            .then(res => {
                console.log(res)
                if(res.success){
                    hasSecret = !!res.data.has_secret
                    needPasswords = !!res.data.need_passwords
                    needSamples = !!res.data.need_samples
                    console.log(hasSecret ? 'Есть секрет' : 'Нет секрета')
                    return hasSecret
                }else{
                    alert(res.error ?? 'Непредвиденная ошибка')
                    return false
                }
            })
            .then(res => {
                showSecretBtn.disabled = (!hasSecret || isEditing)
                return res
            })
    }
    function goCheckSecret(){
        checkSecret()
            .then(hasSecret => {
                showSecretBtn.innerText = hasSecret ? 'Редактировать' : 'Создать секретную записку'
                showSecretBtn.disabled = false
            })
            .then(() => {
                console.log({needSamples, needPasswords, hasSecret})
                if(needPasswords){
                    secretNoteCard.classList.add('d-none')
                    passwordsCard.classList.remove('d-none')
                    renderPasswordsForm()
                    accessLevel.innerText = 'Настройка парольных фраз'
                }else if(needPasswords === false){
                    passwordsCard.classList.add('d-none')
                    if(needSamples){
                        secretNoteCard.classList.add('d-none')
                        accessLevel.innerText = 'Калибровка эталонов для входа по клавиатурному почерку (нужно запускать сессии)'
                    }else{
                        secretNoteCard.classList.remove('d-none')
                        accessLevel.innerText = 'Требуется вход по клавиатурному почерку'
                    }
                }
            })
    }
    goCheckSecret()


    function assertToken(){
        return new Promise(resolve => {
            if(!token){
                // запускаем сессию
                processSession()
                    .then(results => {
                        hideApp()
                        console.log('READY', results)

                        // получаем токен
                        const data = {
                            auth: credentials,
                            flows: results,
                        }
                        fetch('/get-token', {
                            method: 'POST',
                            body: JSON.stringify(data),
                        })
                            .then(res => res.json())
                            .then(res => {
                                if(res.success){
                                    token = res.data.token
                                    console.log('new token: ', token)
                                }else{
                                    alert(res.error ?? 'Непредвиденная ошибка')
                                }
                                resolve()
                            })
                            .catch(err => {
                                console.error(err)
                                resolve()
                            })
                    })
                    .catch(err => {
                        console.error(err)
                        resolve()
                    })
            }else{
                resolve()
            }
        })
    }

    function editSecretNote(){
        showSecretBtn.disabled = true

        assertToken()
            .then(() => {
                if(token){
                    const data = {token}
                    fetch('/user/get-secret', {
                        method: 'POST',
                        body: JSON.stringify(data),
                    })
                        .then(res => res.json())
                        .then(res => {
                            if(res.success){
                                secret.value = res.data
                                secret.disabled = false
                                isEditing = true
                                saveSecretBtn.disabled = false
                                accessLevel.innerText = 'Полный доступ'
                            }else{
                                showSecretBtn.disabled = false
                                alert(res.error ?? 'Непредвиденная ошибка')
                            }
                        })
                }else{
                    showSecretBtn.disabled = false
                }
            })
    }

    function saveSecretNote() {

        const value = secret.value
        if(typeof value !== 'string'){
            alert('Поле должно быть строкой')
            return
        }

        showSecretBtn.disabled = true
        saveSecretBtn.disabled = true

        assertToken()
            .then(() => {
                if(token){
                    const data = {token, value}
                    fetch('/user/set-secret', {
                        method: 'POST',
                        body: JSON.stringify(data),
                    })
                        .then(res => res.json())
                        .then(res => {
                            if(res.success){
                                secret.value = res.data
                                secret.disabled = false
                                isEditing = true
                                saveSecretBtn.disabled = false
                                alert('Секрет сохранён!')
                            }else{
                                showSecretBtn.disabled = false
                                alert(res.error ?? 'Непредвиденная ошибка')
                            }
                        })
                }else{
                    showSecretBtn.disabled = false
                }
            })
    }

    // Попытка просмотреть секретную записку
    showSecretBtn.addEventListener('click', () => editSecretNote())

    // Сохранение секрета
    saveSecretBtn.addEventListener('click', () => saveSecretNote())

})

