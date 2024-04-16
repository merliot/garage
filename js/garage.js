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
		let sensor = this.state.Door.Sensor
		let relay = this.state.Door.Relay

		let range = sensor.Max - sensor.Min
		let percent = 0
		if (range !== 0) {
			percent = Math.round((sensor.Dist - sensor.Min) / range * 100.0 / 5) * 5
			percent = Math.min(100, Math.max(0, percent))
		}

		image.src = "images/door-" + percent + ".png"

		let div = document.getElementById("door")
		div.style.background = (relay.State) ? "cornsilk" : "none"
	}

	setMouse() {
		let div = document.getElementById("door")
		if (this.viewMode === ViewMode.ViewFull) {
			if (this.state.Online) {
				div.onmousedown = () => {
					this.click(true)
				}
				div.onmouseup = () => {
					this.click(false)
				}
			}
		}
	}

	saveClick(msg) {
		this.state.Door.Relay.State = msg.Clicked
		this.setDoorImg()
	}

	savePosition(msg) {
		let sensor = this.state.Door.Sensor
		sensor.Dist = msg.Dist
		sensor.Min = msg.Min
		sensor.Max = msg.Max
		this.setDoorImg()
	}

	click(clicked) {
		let relay = this.state.Door.Relay
		relay.State = clicked
		this.setDoorImg()
		this.webSocket.send(JSON.stringify({Path: "click", Clicked: clicked}))
	}
}
