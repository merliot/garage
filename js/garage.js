import { WebSocketController } from './common.js'

export function run(prefix, url) {
	const garage = new Garage(prefix, url)
}

class Garage extends WebSocketController {

	open() {
		super.open()
		this.showGarage()
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
		for (let i = 0; i < 2; i++) {
			let div = document.getElementById("door" + i)
			let label = document.getElementById("door" + i + "-name")
			var door = this.state.Doors[i]
			if (door.Name === "") {
				div.classList.replace("door", "door-unused")
				label.classList.replace("door-name", "door-name-unused")
				label.textContent = "unasigned"
			} else {
				div.classList.replace("door-unused", "door")
				label.classList.replace("door-name-unused", "door-name")
				label.textContent = door.Name
				div.onmousedown = () => {
					this.click(i, true)
				}
				div.onmouseup = () => {
					this.click(i, false)
				}
			}
			this.setDoorImg(i)
		}
	}

	setDoorImg(index) {
		let image = document.getElementById("door" + index + "-img")
		let door = this.state.Doors[index]

		image.disabled = false

		let range = door.Max - door.Min
		let percent = 0
		if (range !== 0) {
			percent = Math.round((door.Dist - door.Min) / range * 100.0 / 5) * 5
			percent = Math.min(100, Math.max(0, percent))
		}

		image.src = "images/door-" + percent + ".png"

		let div = document.getElementById("door" + index)
		div.style.background = (door.Clicked) ? "cornsilk" : "none"
	}

	saveClick(msg) {
		let door = this.state.Doors[msg.Door]
		door.Clicked = msg.Clicked
		this.setDoorImg(msg.Door)
	}

	savePosition(msg) {
		let door = this.state.Doors[msg.Door]
		door.Dist = msg.Dist
		door.Min = msg.Min
		door.Max = msg.Max
		this.setDoorImg(msg.Door)
	}

	click(index, clicked) {
		var door = this.state.Doors[index]
		door.Clicked = clicked
		this.setDoorImg(index)
		this.webSocket.send(JSON.stringify({Path: "click", Door: index, Clicked: clicked}))
	}
}
