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

async function graphql(body) {
    return await fetch('/api', {
        method: 'POST',
        body
    }).then(async r => await r.json())
}

let commands = []

window.onload = () => {
    document.getElementById("main").setAttribute('style', 'background-image: url(/images/Space_parallax.png)')
    graphql(`query{ commands { name description } }`).then(r => {
        commands = r.data.commands
        updateView()
    })
}

function updateView() {
    let table = document.getElementById('commands-table')
    for(let i = 0; i < commands.length; i++) {
        let cmd = commands[i]
        let tr = document.createElement("tr")
        let td_name = document.createElement("td")
        let td_desc = document.createElement("td")
        td_name.innerHTML = cmd.name
        td_desc.innerHTML = cmd.description
        
        tr.appendChild(td_name)
        tr.appendChild(td_desc)
        
        table.appendChild(tr)
    }
}
