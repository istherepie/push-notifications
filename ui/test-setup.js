import Jasmine from 'jasmine'

const jasmine = new Jasmine()

const config = {
	spec_dir: 'tests',
    spec_files: [
		"**/*[sS]pec.js"
    ],
    helpers: [
        "helpers/**/*.js"
	],
	stopSpecOnExpectationFailure: false,
	random: true
}

jasmine.loadConfig(config)
jasmine.execute()
