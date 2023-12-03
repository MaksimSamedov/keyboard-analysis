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

    // Основные элементы интерфейса
    const start = document.getElementById('start')
    const stop = document.getElementById('stop')
    const remain = document.getElementById('remain')
    const app = document.getElementById('app')

    // Форма ввода паролей
    const result = document.getElementById('result')
    const currentPhrase = document.querySelector('label[for=result]')


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
                    .then(console.log)
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



    // История

    const history = document.getElementById('history')
    let selected = null
    function reloadHistory() {
        history.innerHTML = ''
        if(!selected){
            fetch('/history', {
                method: 'POST',
                body: JSON.stringify(credentials),
            }).then(res => res.json())
                .then(res => {
                    console.log(res)
                    res.forEach(flow => {
                        const el = document.createElement('div')
                        el.innerText = flow.phrase
                        history.append(el)
                        el.addEventListener('click', e => {
                            selected = flow.id
                            reloadHistory()
                        })
                    })
                })
        }else{
            fetch('/history/' + selected, {
                method: 'POST',
                body: JSON.stringify(credentials),
            })
                .then(res => res.json())
                .then(res => {
                    console.log(res)
                })
        }
    }
    reloadHistory()


    // Секретная записка
    const secret = document.getElementById('secret-note')
    const showSecretBtn = document.getElementById('secret-show')
    const saveSecretBtn = document.getElementById('secret-save')

    let isEditing = false
    let hasSecret = null
    function checkSecret(){
        showSecretBtn.disabled = true
        return fetch('/user/has-secret', {
            method: 'POST',
            body: JSON.stringify(credentials),
        }).then(res => res.json())
            .then(res => {
                console.log(res)
                if(res.success){
                    hasSecret = !!res.data
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
    checkSecret()
        .then(hasSecret => {
            showSecretBtn.innerText = hasSecret ? 'Редактировать' : 'Создать секретную записку'
            showSecretBtn.disabled = false
        })


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

