'use strict';

const DEFAULT_ALPHABET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'

document.ready(() => {

    // Основные элементы интерфейса
    const start = document.getElementById('start')
    const stop = document.getElementById('stop')
    const remain = document.getElementById('remain')
    const app = document.getElementById('app')
    const launcher = document.getElementById('launcher')

    // Форма запуска сессии
    const phrases = document.getElementById('phrases')
    const length = document.getElementById('length')
    const alphabet = document.getElementById('alphabet')
    alphabet.value = DEFAULT_ALPHABET

    // Форма ввода паролей
    const result = document.getElementById('result')
    const currentPhrase = document.querySelector('label[for=result]')


    // Отобразить форму запуска
    const showLauncher = () => {
        app.classList.add('d-none')
        launcher.classList.remove('d-none')

        start.classList.remove('d-none')
        stop.classList.add('d-none')
    }
    // Отобразить форму ввода
    const showApp = () => {
        app.classList.remove('d-none')
        launcher.classList.add('d-none')

        start.classList.add('d-none')
        stop.classList.remove('d-none')
    }



    // События

    let session = null
    start.addEventListener('click', () => {
        session = null
        showLauncher()
        if(phrases.value < 1 || length.value < 2 || alphabet.value.length < 10){
            return
        }
        session = newSession({
            phrases: randomPhrases(length.value * 1, phrases.value * 1, alphabet.value),
            callback: (results) => {
                showLauncher()
                console.log('READY', results)
                fetch('/process', {
                    method: 'POST',
                    body: JSON.stringify(results),
                })
                    .then(res => res.json())
                    .then(console.log)
            },
        })
        showApp()
        session.next()
    })

    stop.addEventListener('click', () => {
        session = null
        showLauncher()
    })

    result.addEventListener('keydown', e => session && session.track(e))
    result.addEventListener('keyup', e => session && session.track(e))



    const isPhraseReady = () => {
        return session.current === result.value
    }



    // Создание сессии
    function newSession({phrases, callback}){
        let row = []

        const sess = {
            current: '',
            result: [],
            source: phrases,

            track: (keyboardEvent) => {
                if(keyboardEvent.repeat){
                    return
                }

                row.push({
                    key: keyboardEvent.key,
                    time: Math.round(keyboardEvent.timeStamp),
                    up: (keyboardEvent.type === 'keyup'),
                })

                if(isPhraseReady()){
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

    function randomPhrases(length, count, characters){
        const phrases = []
        for(let i = 0; i < count; i++){
            phrases.push(randomString(length, characters))
        }
        return phrases
    }

    function randomString(length, characters){
        let result = ''
        const charactersLength = characters.length
        let counter = 0
        while (counter < length) {
            result += characters.charAt(Math.floor(Math.random() * charactersLength))
            counter += 1
        }
        return result
    }



    // История

    const history = document.getElementById('history')
    let selected = null
    function reloadHistory() {
        history.innerHTML = ''
        if(!selected){
            fetch('/history').then(res => res.json())
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
            fetch('/history/' + selected)
                .then(res => res.json())
                .then(res => {

                })
        }
    }
    reloadHistory()

})

