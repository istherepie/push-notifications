/**
 * Copyright (c) 2020
 *
 * Boilerplate: MainJS file template
 *
 * @summary an awesome template
 * @author Steffen Park <dev@istherepie.com>
 */

import 'modern-css-reset'
import '$base/style.css'

class App {
	constructor(elementID) {
		this.root = document.getElementById(elementID)
		this.input = document.getElementById("notification-input")
		this.area = document.getElementById("notification-area")
		this.socket = new EventSource("/api/notifications")
	}

	createNotification(message) {
		// Create notification container
		let notification = document.createElement("div")
		notification.classList = "notification fade-in"

		// Create title
		let title = document.createElement("h4")
		title.innerHTML = "Notification"
		title.classList = "notification-title"

		// Create body
		let body = document.createElement("div")
		body.innerHTML = message

		// Append title and body
		notification.append(title, body)
		return notification
	}

	notify(message) {
		let notification = this.createNotification(message)
		this.area.append(notification)

		setTimeout(() => {
			notification.classList = "notification fade-out"
			notification.addEventListener("animationend", event => {
				notification.remove()
			})
			
		}, 6000)
	}

	fire(event) {
		event.preventDefault()
		console.log(event)

		this.request()
	}

	request() {
		var req = new XMLHttpRequest();

		req.addEventListener("load", event => {
			console.log(event)
		});

		req.onreadystatechange = event => {
			console.log(event)
		}

		let payload = {
			message: this.input.value
		}

		let postData = JSON.stringify(payload)

		req.open("POST", "/api/message");
		req.send(postData);
	}

	Run() {
		document.body.append(this.area)

		this.socket.onmessage = (event) => {this.notify(event.data)}

		let button = document.getElementById("submit")
		button.addEventListener("click", event => this.fire(event));
	}
}

// Init the app!
const app = new App("root")
window.onload = () => { app.Run() }
