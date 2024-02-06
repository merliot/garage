import { WebSocketController, ViewMode } from './common.js'

export function run(prefix, url, viewMode) {
	const garage = new Garage(prefix, url, viewMode)
}

class Garage extends WebSocketController {

	open() {
		super.open()
		if (this.state.DeployParams !== "") {
			this.showGarage()
		}
	}

	handle(msg) {
		switch(msg.Path) {
		case "click":
			this.saveClick(msg)
			break
		case "position":
			this.savePosition(msg)
			break
		}
	}

	showGarage() {
		let div = document.getElementById("door")

		if (this.state.Door.Name === "") {
			div.classList.replace("visibleFlex", "hidden")
		} else {
			this.setDoorName()
			this.setDoorImg()
			this.setMouse()
			div.classList.replace("hidden", "visibleFlex")
		}
	}

	setDoorName() {
		let label = document.getElementById("door-name")
		label.textContent = this.state.Door.Name
	}

	setDoorImg() {
		let image = document.getElementById("door-img")
		let door = this.state.Door

		let range = door.Max - door.Min
		let percent = 0
		if (range !== 0) {
			percent = Math.round((door.Dist - door.Min) / range * 100.0 / 5) * 5
			percent = Math.min(100, Math.max(0, percent))
		}

		image.src = "images/door-" + percent + ".png"

		let div = document.getElementById("door")
		div.style.background = (door.Clicked) ? "cornsilk" : "none"
	}

	setMouse() {
		let div = document.getElementById("door")
		if (this.viewMode === ViewMode.ViewFull) {
			div.onmousedown = () => {
				this.click(true)
			}
			div.onmouseup = () => {
				this.click(false)
			}
		}
	}

	saveClick(msg) {
		this.state.Door.Clicked = msg.Clicked
		this.setDoorImg()
	}

	savePosition(msg) {
		let door = this.state.Door
		door.Dist = msg.Dist
		door.Min = msg.Min
		door.Max = msg.Max
		this.setDoorImg()
	}

	click(clicked) {
		let door = this.state.Door
		door.Clicked = clicked
		this.setDoorImg()
		this.webSocket.send(JSON.stringify({Path: "click", Clicked: clicked}))
	}
}
