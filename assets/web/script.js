
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
    const speed = 10
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

const force_auth = false
const dummy = {
    username: "Sweetheart",
    discriminator: "1857",
    picture: "/images/sweetheart.png",
    profile: {
        about: "I'm absolutely amazing!",
        description: "Hi, i'm **Sweetheart**!\nAnd i'm absolutely amazing *obviously*!",
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
        timezone: new Date().toTimeString().substring(0, 5),
        country: "Netherlands",
        
        gender: "???",
        pronouns: "???/???/???",
        sexuality: "???"
    }
}

const dummy_countries = [
    {
        name: "Netherlands",
        flag: "🇳🇱"
    },
    {
        name: "Germany",
        flag: "🇩🇪"
    }, 
    {
        name: "China",
        flag: "🇨🇳"
    }
]

let user = {}
let countries = []

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
        countries = dummy_countries
        authenticated()
        return
    }
    if(localStorage.discord_session) {
        console.log('authenticated')
        graphql(`{identity(session: "${localStorage.discord_session}") { username discriminator picture profile{ about description favorite_color socials { name handle } timezone country gender pronouns sexuality } } countries { name flag }}`).then(r => {
            user = r.data.identity
            countries = r.data.countries
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
    document.getElementById('login').style.display = 'none'
    document.getElementById('inputs').style.display = 'grid'
    let picker = Pickr.create({
        el: '#preview-color',
        theme: 'nano',
        default: `#${user.profile.favorite_color.toString(16)}`,
        components: {
            preview: true,
            hue: true,
            interaction: {
                hex: true,
                input: true,
                save: true,
            }
        }
    })
    picker.on('save', (color, instance) => {
        user.profile.favorite_color = parseInt(color.toHEXA().join(""), 16)
        updatePreview()
    })
    
    countries.sort((a, b) => a.name.localeCompare(b.name))
    
    let celem = document.getElementById('preview-country')
    for(let i = 0; i < countries.length; i++) {
        let c = document.createElement("option")
        c.innerText = `${countries[i].name} ${countries[i].flag}`
        c.value = countries[i].name
        celem.appendChild(c)
    }
    celem.value = user.profile.country
    
    updatePreview()
}

/*
"__NAME:__ HANDLE"

<li>
	<div>
	    {MD}
	</div>
	<button onclick="removeSocial({INDEX})">
		<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="20" height="20">
		    <line x1="3" x2="17" y1="10" y2="10" stroke="rgb(220, 221, 222)" stroke-width="2" stroke-linecap="round"></line>
		</svg>
	</button>
</li>
*/

function updatePreview() {
    user.profile.timezone = new Date().toTimeString().substring(0, 5)
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
            if(elem.getAttribute('raw') === 'true') {
                elem.innerHTML = prv
                continue
            }
            if(elem.tagName.toLowerCase() === "input") {
                elem.value = prv
                continue
            }
            elem.innerHTML = discordMarkdown.toHTML(prv)
        }
    }
    
    const elem = document.getElementById("preview-socials")
    const socPrev = document.getElementById("socials-preview")
    socPrev.innerHTML = ''
    let social = ''
    if(!user.profile.socials) return
    if(user.profile.socials.length == 0) {
        elem.innerHTML = 'no socials.'
        return
    }
    for(let i = 0; i < user.profile.socials.length; i++) {
        let s = user.profile.socials[i]
        social += `> __${s.name}:__ ${s.handle}\n`
        let li = document.createElement("li")
        let button = document.createElement("button")
        let div = document.createElement("div")
        div.innerHTML = discordMarkdown.toHTML(`__${s.name}:__ ${s.handle}`)
        button.innerHTML = `<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="20" height="20"><line x1="3" x2="17" y1="10" y2="10" stroke="rgb(220, 221, 222)" stroke-width="2" stroke-linecap="round"></line></svg>`
        button.setAttribute('onclick', `deleteSocial(${i})`)
        li.appendChild(div)
        li.appendChild(button)
        socPrev.appendChild(li)
    }
    
    elem.innerHTML = discordMarkdown.toHTML(social)
}

function getDetail(path) {
    let dir = path.split(".")
    let root = user
    for(let i = 0; i < dir.length; i++) {
        let v = root[dir[i]]
        if(v != undefined && v != null) {
            root = root[dir[i]]
        }else{
            break
        }
    }
    if(path.endsWith("country")) {
        root += ` ${countries.find((e) => e.name === root).flag}`
    }
    return root
}

function reset() {
    authenticated()
}

function updateValue(el, path) {
    let up = {}
    up[path] = el.value ? el.value : ''
    user.profile = Object.assign(user.profile, up)
    updatePreview()
}

function save() {
    graphql(JSON.stringify({
        query: `mutation($about: String! $description: String! $favorite_color: Int! $timezone: Int! $gender: String! $pronouns: String! $sexuality: String! $country: String!) {
                profile(session: "${localStorage.discord_session}" 
                    about: $about  
                    description: $description 
                    favorite_color: $favorite_color 
                    socials: $socials
                    timezone: $timezone 
                    gender: $gender 
                    pronouns: $pronouns 
                    sexuality: $sexuality 
                    country: $country
                )
            }`,
        variables: Object.assign(user.profile, {
            timezone: new Date().getTimezoneOffset()
        })
    })).then(r => {
        if(r.data.profile) {
            notify("Saved profile!")
        }else{
            notify("Failed to save profile! A warning has been sent to the developer about the issue!")
        }
    })
}

let notif = null

function notify(message) {
    document.getElementById("notification").innerHTML = message
    document.getElementById("notification").style.opacity = 1
    if(notif) clearTimeout(notif)
    
    notif = setTimeout(() => {
        document.getElementById("notification").style.opacity = 0
    }, 5000)
}

function addSocial() {
    let name = document.getElementById("preview-social-name")
    let handle = document.getElementById("preview-social-handle")
    if(!(name.value.length > 0 && handle.value.length > 0)) return
    user.profile.socials.push({
        name: name.value.trim(),
        handle: handle.value.trim()
    })
    name.value = ""
    handle.value = ""
    updatePreview()
}

function deleteSocial(index) {
    delete user.profile.socials[index]
    let n = []
    for(let i = 0; i < user.profile.socials.length; i++) {
        if(user.profile.socials[i]) n.push(user.profile.socials[i])
    }
    user.profile.socials = n
    updatePreview()
}
