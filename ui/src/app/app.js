export class App {
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