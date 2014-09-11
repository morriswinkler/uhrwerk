Fab Lab Booking HTML Interface
==============================

This is [Grunt](http://gruntjs.com) powered HTML interface building sandbox. Read this README before doing any harm to it.

1. Clone the repo
2. Go to `html/` directory (you are here already)
3. Install Grunt locally with `npm install` (the Node Package Manager will install things based on contents in `package.json`)
4. Install [Bower](http://bower.io) with `bower install` (packages are defined in `bower.json`, they will be installed in `src/libs`)
5. Use `grunt dev` while developing - it watches for changes in the src directory and auto-reloads the src/index.html file if any change occurs
6. Use `grunt dist` to compile, minify js and css files, replace multiple resource instances in `src/index.html` with single ones and copy all production-ready files to the `dist/` folder.