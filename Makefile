UGLIFY = ./node_modules/uglify-js/bin/uglifyjs
CLEANCSS = ./node_modules/clean-css/bin/cleancss

cloc:
	cloc . --exclude-dir=app/bower_components,app/scripts/.module-cache

static:
	$(UGLIFY) \
		app/bower_components/react/react.min.js \
		app/bower_components/react-router/build/global/ReactRouter.min.js \
		app/bower_components/lodash/lodash.min.js \
		app/js/* \
		app/jsx/build/charts/* \
		app/jsx/build/app.js \
		> build/app.js
	cat \
		app/bower_components/normalize.css/normalize.css \
		app/css/* \
		| $(CLEANCSS) -o build/app.css
