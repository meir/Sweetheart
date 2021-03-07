
window.bg_towards = 0;
window.bg_last = 0;

(() => {
    let lastTime = 0
    const second = 1000
    const frames = 60
    const jumps = 20
    let f = (t) => {
        if((t - lastTime) >= second/frames) {
            let x = window.bg_last
            x -= (x-window.bg_towards)/jumps
            window.bg_last = x
            
            document.body.style.backgroundPositionX = `${x}%`
            lastTime = t
        }
        requestAnimationFrame(f)
    }
    f()
})()

window.onmousemove = (e) => {
    const speed = 20
    let x = e.clientX / speed
    window.bg_towards = x
}

//--

async function graphql(body) {
    return await fetch('/api', {
        method: 'POST',
        body
    }).then(async r => await r.json())
}

window.onload = () => {
    const query = `{settings{oauth invite}}`
    graphql(query).then(r => {
        document.getElementById('login').setAttribute('a', r.data.settings.oauth)
        document.getElementById('invite').setAttribute('a', r.data.settings.invite)
    })
}
