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
	}

	Hello() {
		let title = document.createElement("h1")
		title.innerHTML = "Hello World"
		this.root.append(title)
	}
}

window.onload = () => {
	new App("root").Hello()
}