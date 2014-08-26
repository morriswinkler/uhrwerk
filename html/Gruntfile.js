'use strict';

module.exports = function (grunt) {

  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),
		
		// package options
		jshint: {
		  options: {
		    jshintrc: '.jshintrc'
		  },
		  all: [
		    'Gruntfile.js',
		    'src/assets/js/*.js'
		  ]
		},
		less: {
			dev: {
				files: {
          'src/assets/css/main.css': ['src/assets/less/main.less']
        },
        options: {
        	compress: false
        }
			},
      dist: {
        files: {
          'dist/assets/css/main.min.css': [
          	'src/libs/bootstrap/less/bootstrap.less',
          	'src/assets/less/main.less'
          ]
        },
        options: {
        	compress: true,
        	cleancss: true
        }
      }
    },
		concat: {
		  dist: {
		    src: [
		      'src/libs/jquery/dist/jquery.js',
		      'src/libs/bootstrap/dist/js/bootstrap.js',
		      'src/assets/js/main.js'
		    ],
		    dest: 'tmp/app.js'
		  },
		  modernizr: {
		    src: [
		      'src/libs/modernizr/modernizr.js'
		    ],
		    dest: 'tmp/modernizr.js'
		  }
		},
		uglify: {
		  dist: {
		    files: {
		      'dist/assets/js/modernizr.min.js' : 'tmp/modernizr.js',
		      'dist/assets/js/app.min.js' : 'tmp/app.js'
		    }
		  }
		},
		processhtml: {
			dist: {
				options: {
					process: true
				},
				files: {
					'dist/index.html': ['src/index.html']
				}
			}
		},
		copy: {
			dist: {
				files: [
					{
						expand: true, 
						src: ['src/libs/bootstrap/dist/fonts/*'], 
						dest: 'dist/assets/fonts/', 
						flatten: true,
						filter: 'isFile'
					}
				]
			}
		},
		clean: {
		  dist: [
		    'tmp/**'
		  ]
		},
		watch: {
			js: {
			  files: [
			    'src/assets/js/*.js',
			    'src/*.html',
			    'src/assets/less/*.less'
			  ],
			  tasks: ['jshint', 'less:dev'],
			  options: {
			    livereload: true,
			    atBegin: true
			  }
			},
		}
  });

	// Load tasks
	grunt.loadNpmTasks('grunt-contrib-clean');
	grunt.loadNpmTasks('grunt-contrib-jshint');
	grunt.loadNpmTasks('grunt-contrib-concat');
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-notify');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-processhtml');
	grunt.loadNpmTasks('grunt-contrib-less');
	grunt.loadNpmTasks('grunt-contrib-copy');


	// Register tasks
  grunt.registerTask('default', ['jshint', 'less', 'concat', 'uglify', 'processhtml', 'copy', 'clean']);
  grunt.registerTask('dist', ['jshint', 'less:dist', 'concat', 'uglify', 'processhtml', 'copy', 'clean']);
  grunt.registerTask('dev', ['watch']);

};