export class EventHandler {
	constructor(eventSource) {
		this.es = eventSource
		this.lastUpdated = Date.now()
		this.monitorConnection()
	}

	subscribe(handler) {
		this.es.addEventListener("message", event => {
			handler("message", event.data)
		})

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

export class MessageHandler {
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

export class NotificationHandler {
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
