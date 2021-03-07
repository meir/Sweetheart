
window.bg_towards = 0;
window.bg_last = 0;

(() => {
    let lastTime = 0
    const second = 1000
    const frames = 60
    const jumps = 20
    let f = (t) => {
        if(window.stop_loop === true) return
        if((t - lastTime) >= second/frames) {
            let x = window.bg_last
            let change = (x-window.bg_towards)/jumps
            if(change >= .02 || change <= -.02) {
                x -= change
                window.bg_last = x
            
                document.documentElement.style.setProperty('--bg-offset', `${x}%`)
            }
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
    document.getElementById("main").setAttribute('style', 'background-image: url(/images/Space_parallax.png)')
    const query = `{settings{oauth invite}}`
    graphql(query).then(r => {
        document.getElementById('login').setAttribute('href', r.data.settings.oauth)
        document.getElementById('invite').setAttribute('href', r.data.settings.invite)
    })
    
    if(localStorage.discord_session) {
        console.log('authenticated')
    }else{
        const urlParams = new URLSearchParams(window.location.search);
        const discordCode = urlParams.get('code');
        
        if(discordCode) {
            const code_query = `{auth(code: "${discordCode}")}`
        
            graphql(code_query).then(r => {
                localStorage.discord_session = r.data.auth
                window.location.href = "/"
            })
        }
    }
}
