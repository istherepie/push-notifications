/**
 * Copyright (c) 2020
 *
 * Boilerplate: MainJS file template
 *
 * @summary an awesome template
 * @author Steffen Park <dev@istherepie.com>
 */

// CSS Reset
import 'modern-css-reset'

// Custom
import '$base/style.css'

// Fontawesome
import { library, dom } from '@fortawesome/fontawesome-svg-core'
import { faRocket, faArrowRight } from '@fortawesome/free-solid-svg-icons'
import { faEnvelope, faCopyright } from '@fortawesome/free-regular-svg-icons'
import { faGithub } from '@fortawesome/free-brands-svg-icons'


class EventHandler {
	constructor(eventSource) {
		this.es = eventSource
		this.lastUpdated = Date.now()
		this.monitorConnection()
	}

	subscribe(handler) {
		this.es.onmessage = (event) => {
			handler("message", event.data)
		}

		this.es.addEventListener("service", event => {
			handler("service", event.data)
		})
	}

	monitorConnection() {
		this.es.addEventListener("heartbeat", event => {
			let now = Date.now() 

			let diff = now - this.lastUpdated

			if ( diff > 15000 ) {
				alert("Connection has been lost!")
				return
			}

			this.lastUpdated = now
		})
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
		return new Promise((resolve, reject) => {

			http.upload.addEventListener("load", complete => {
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

	createNotification(msgValue) {
		// Create pContainer
		let p = document.createElement("p")
		p.innerHTML = msgValue

		let div = document.createElement("div")
		div.append(p)

		// Create notification container
		let notification = document.createElement("div")
		

		// Append title and body
		notification.append(div)
		return notification
	}

	fire(msgType, msgValue) {
		let notification = this.createNotification(msgValue)
		notification.classList = `notification ${msgType}`
		this.root.append(notification)

		setTimeout(() => {
			notification.remove()
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
		
		this.messageHandler.publish(message)
		.then(result => {
			// do nothing
		})
		.catch(err => {
			console.log(err)
			this.notificationHandler.fire("service", "Connection to the service failed")
		})
	}


	Run() {
		let trigger = document.getElementById("submit")

		trigger.addEventListener("click", event => {
			event.preventDefault()
			this.example()
		});

		this.eventHandler.subscribe((msgType, msgValue) => {
			this.notificationHandler.fire(msgType, msgValue)
		})

		const urlTarget = document.getElementById("baseurl")
		urlTarget.innerHTML = new URL("/message", window.location)

	}
}

// MAIN()
window.onload = () => {

	// Icons
	library.add(faArrowRight, faRocket, faGithub, faEnvelope, faCopyright)
	dom.watch()

	// Setup Handlers
	const messageHandler = new MessageHandler("/api/message")
	const notificationHandler = new NotificationHandler("notification-area")

	// EventSource API
	const es = new EventSource("/api/notifications")
	const eventHandler = new EventHandler(es)

	// Init App
	const app = new App(eventHandler, messageHandler, notificationHandler)
	app.Run() 

}
