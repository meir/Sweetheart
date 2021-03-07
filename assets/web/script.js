
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
