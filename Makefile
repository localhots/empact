JSX = ./node_modules/react-tools/bin/jsx
UGLIFY = ./node_modules/uglify-js/bin/uglifyjs
CLEANCSS = ./node_modules/clean-css/bin/cleancss

cloc:
	cloc . --exclude-dir=node_modules,app/bower_components,app/jsx/build/.module_cache

dist:
	# Compiling JSX
	mkdir -p build
	rm -rf app/jsx/build
	$(JSX) --extension jsx app/jsx app/jsx/build

	# Compressing JS
	cat app/index.html \
		| grep script \
		| grep -v bower \
		| cut -d '"' -f 2 \
		| sed -e 's/^/app/' \
		| ( echo app/bower_components/react/react.min.js && \
			echo app/bower_components/react-router/build/global/ReactRouter.min.js && \
			echo app/bower_components/lodash/lodash.min.js && \
			cat ) \
		| xargs cat \
		| $(UGLIFY) > build/app.js

	# Compressing CSS
	cat app/index.html \
		| grep stylesheet \
		| grep -v opensans \
		| cut -d '"' -f 4 \
		| sed -e 's/^/app/' \
		| xargs cat \
		| $(CLEANCSS) -o build/app.css

	# Replacing all link tags with a single one
	# Replacing all script tags with a single one
	cat app/index.html \
		| perl -pe 's/<link.*>\n/STYLE_TAG/g' \
		| perl -pe 's/(STYLE_TAG\s*)+/<link rel="stylesheet" href="\/app.css">\n        <link rel="stylesheet" href="http:\/\/fonts.googleapis.com\/css?family=Open\+Sans:400,300,600">\n    /g' \
		| perl -pe 's/<script.*>\n/SCRIPT_TAG/g' \
		| perl -pe 's/(SCRIPT_TAG\s*)+/<script type="text\/javascript" src="\/app.js"><\/script>\n/g' \
		> build/index.html
