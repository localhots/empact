JSX = ./node_modules/react-tools/bin/jsx
JSMIN = ./node_modules/uglify-js/bin/uglifyjs
CSSMIN = ./node_modules/clean-css/bin/cleancss

cloc:
	cloc . --exclude-dir=node_modules,app/bower_components,app/jsx/build/.module_cache

dist:
	# Compiling JSX
	rm -rf app/jsx/build
	$(JSX) --extension jsx app/jsx app/jsx/build

	# Compressing JS
	#
	# Get all script tags
	# Ignore Bower deps, include minified versions later
	# Splitting string by " and taking second field (URI)
	# Adding "app" prefix
	# Prepending dependencies to the list
	# Concatenation
	# Minification
	cat app/app.html \
	| grep script \
	| grep -v bower \
	| cut -d '"' -f 2 \
	| sed -e 's/^/app/' \
	| ( echo app/bower_components/react/react.min.js && \
		echo app/bower_components/react-router/build/global/ReactRouter.min.js && \
		echo app/bower_components/lodash/lodash.min.js && \
		cat ) \
	| xargs cat \
	| $(JSMIN) > build/app.js

	# Compressing CSS
	#
	# Get all style tags
	# Ignore local Open Sans stylesheet
	# Splitting string by " and taking fourth field (URI)
	# Adding "app" prefix
	# Concatenation
	# Minification
	cat app/app.html \
	| grep stylesheet \
	| grep -v fonts.googleapis.com \
	| grep -v opensans \
	| cut -d '"' -f 4 \
	| sed -e 's/^/app/' \
	| xargs cat \
	| $(CSSMIN) -o build/app.css
