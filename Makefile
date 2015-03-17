UGLIFY = ./node_modules/uglify-js/bin/uglifyjs
CLEANCSS = ./node_modules/clean-css/bin/cleancss

cloc:
	cloc . --exclude-dir=app/bower_components,app/scripts/.module-cache

static:
	$(UGLIFY) \
		app/bower_components/react/react.min.js \
		app/bower_components/react-router/build/global/ReactRouter.js \
		app/bower_components/lodash/lodash.js \
		app/js/* \
		app/jsx/build/*/* \
		> build/app.js
	cat app/css/* | $(CLEANCSS) -o build/app.css
