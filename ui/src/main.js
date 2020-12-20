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
import "@fortawesome/fontawesome-free/js/all"


class EventHandler {
	constructor(eventURI) {
		this.es = new EventSource(eventURI)
	}

	// The handler must accept a message (string)
	subscribe(handler) {
		this.es.onmessage = (event) => {
			handler(event.data)
		}
	}
}

class MessageHandler {
	constructor(messageURI) {
		this.uri = messageURI
	}

	publish(message) {
		if (typeof message != "string") {
			throw new Error("Message must be passed as a string")
		}

		let payload = JSON.stringify({
			message: message
		});

		return this.request(payload)
	}

	request(payload) {
		var http = new XMLHttpRequest()
		return new Promise((reject, resolve) => {

			http.upload.addEventListener("load", complete => {
				console.log(complete)
				return resolve(complete)
			});

			http.upload.addEventListener("error", err => {
				return reject(err)
			});
		
			http.open("POST", this.uri)
			http.send(payload);
		})

	}
}

class NotificationHandler {
	constructor(elementID) {
		this.root = document.getElementById(elementID)
	}

	createNotification(message) {
		// Create iconContainer
		let iconContainer = document.createElement("div")
		iconContainer.classList = "icon"

		let icon = document.createElement("i")
		icon.classList = "fas fa-comment-alt fa-2x white"

		iconContainer.append(icon)

		// Create pContainer
		let p = document.createElement("p")
		p.classList = "notification-text"
		p.innerHTML = message

		let pContainer = document.createElement("div")
		pContainer.append(p)

		// Create notification container
		let notification = document.createElement("div")
		notification.classList = "notification fade-in"

		// Append title and body
		notification.append(iconContainer, pContainer)
		return notification
	}

	fire(message) {
		let notification = this.createNotification(message)
		this.root.append(notification)

		setTimeout(() => {
			notification.classList = "notification fade-out"
			notification.addEventListener("animationend", event => {
				notification.remove()
			})
			
		}, 6000)
	}
}

class App {
	constructor(eventHandler, messageHandler, notificationHandler) {
		this.eventHandler = eventHandler
		this.messageHandler = messageHandler
		this.notificationHandler = notificationHandler
	}

	example() {
		let message = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
		return this.createEvent(message)
	}

	createEvent(message) {
		this.messageHandler.publish(message)
		.then(result => {
			console.log(result)
		})
	}

	Run() {
		let trigger = document.getElementById("submit")

		trigger.addEventListener("click", event => {
			event.preventDefault()
			this.example()
		});

		this.eventHandler.subscribe(message => {
			this.notificationHandler.fire(message)
		})

	}
}

// Init
const eventHandler = new EventHandler("/api/notifications")
const messageHandler = new MessageHandler("/api/message")
const notificationHandler = new NotificationHandler("notification-area")

const app = new App(eventHandler, messageHandler, notificationHandler)

// MAIN()
window.onload = () => {
	app.Run() 
}
