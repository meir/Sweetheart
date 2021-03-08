
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

const force_auth = true
const dummy = {
    username: "Sweetheart",
    discriminator: "1857",
    picture: "/images/sweetheart.png",
    profile: {
        about: "I'm absolutely amazing!",
        description: "Hi, i'm **Sweetheart** and i'm absolutely amazing *obviously*!",
        favorite_color: 0xffffff,
        socials: [
            {
                name: 'Discord',
                handle: 'Sweetheart#1857'
            },
            {
                name: 'Website',
                handle: 'https://sweetheart.flamingo.dev/'
            }
        ],
        timezone: "CET",
        country: "Netherlands",
        
        gender: "???",
        pronouns: "???/???/???",
        sexuality: "???"
    }
}

let user = {}

window.onload = () => {
    document.getElementById("main").setAttribute('style', 'background-image: url(/images/Space_parallax.png)')
    const query = `{settings{oauth invite}}`
    graphql(query).then(r => {
        document.getElementById('login').setAttribute('href', r.data.settings.oauth)
        document.getElementById('invite').setAttribute('href', r.data.settings.invite)
    })
    
    if(force_auth) {
        console.log('authenticated')
        user = dummy
        authenticated()
        return
    }
    if(localStorage.discord_session) {
        console.log('authenticated')
        graphql(`{identity(session: "${localStorage.discord_session}") { username discriminator picture profile{ about description favorite_color timezone country gender pronouns sexuality } }}`).then(r => {
            user = r.data.identity
            authenticated()
        })
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

function authenticated() {
    document.getElementById('no-preview').style.display = 'none'
    document.getElementById('embed').style.display = 'grid'
    updatePreview()
}

function updatePreview() {
    const elems = document.getElementsByClassName("preview")
    for(let i = 0; i < elems.length; i++) {
        let elem = elems[i]
        if(!elem) continue
        if(elem.getAttribute('preview-attr')) {
            let val = getDetail(elem.getAttribute("preview"))
            val = typeof val == "number" ? `#${val.toString(16)}` : val
            elem.style[elem.getAttribute('preview-attr')] = val
        }else{
            if(elem.tagName.toLowerCase() === "img") {
                elem.src = getDetail(elem.getAttribute("preview"))
                continue
            }
            let prv = getDetail(elem.getAttribute("preview"))
            if(elem.getAttribute('preview-tmpl')) {
                prv = elem.getAttribute('preview-tmpl').replace('%v', prv)
            }
            elem.innerHTML = discordMarkdown.toHTML(prv)
        }
    }
    
    const elem = document.getElementById("preview-socials")
    let social = ''
    if(!user.profile.socials) return
    for(let i = 0; i < user.profile.socials.length; i++) {
        let s = user.profile.socials[i]
        social += `> __${s.name}:__ ${s.handle}\n`
    }
    elem.innerHTML = discordMarkdown.toHTML(social)
    
}

function getDetail(path) {
    let dir = path.split(".")
    let root = user
    for(let i = 0; i < dir.length; i++) {
        if(root[dir[i]]) {
            root = root[dir[i]]
        }else{
            break
        }
    }
    return root
}
