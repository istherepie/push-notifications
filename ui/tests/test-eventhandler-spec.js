import { EventHandler, NotificationHandler } from "../src/app/handlers"

class MockEventSource {

	constructor() {
		this.handlers = {}
	}

	addEventListener(name, handler) {
		this.handlers[name] = handler

	}

	inject(message) {

		let event = {
			data: message
		}

		let handler = this.handlers["message"]
		handler(event)
	}

}

describe("Test EventHandler", () => {

	it("Should receive and forward a message", () => {

		const mockEs = new MockEventSource()
		const eventHandler = new EventHandler(mockEs) 

		let testMessage = "this is a test message"

		let result = ""

		eventHandler.subscribe((msgType, msgValue) => {
			result = msgValue
		})

		mockEs.inject(testMessage)

		expect(result).toBe(testMessage)
	});
});