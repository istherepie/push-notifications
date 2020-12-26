/**
 * Copyright (c) 2020
 *
 * Boilerplate: MainJS file template
 *
 * @summary an awesome template
 * @author Steffen Park <dev@istherepie.com>
 */

// CSS Reset
import "modern-css-reset"

// Custom
import "$base/style.css"

// Fontawesome
import { library, dom } from "@fortawesome/fontawesome-svg-core"
import { faRocket, faArrowRight } from "@fortawesome/free-solid-svg-icons"
import { faEnvelope, faCopyright } from "@fortawesome/free-regular-svg-icons"
import { faGithub } from "@fortawesome/free-brands-svg-icons"

// Application
import { App } from "$app/app"
import { MessageHandler, NotificationHandler, EventHandler } from "$app/handlers"

// MAIN()
window.onload = () => {
	// Icons
	library.add(faArrowRight, faRocket, faGithub, faEnvelope, faCopyright)
	dom.watch()

	// EventSource API
	const es = new EventSource("/event/notifications")

	// Setup Handlers
	const eventHandler = new EventHandler(es)
	const messageHandler = new MessageHandler("/event/message")
	const notificationHandler = new NotificationHandler("notification-area")

	// Init App
	const app = new App(eventHandler, messageHandler, notificationHandler)
	app.Run() 
}
