#!/usr/bin/env node

( async () => {
	const process = require( "process" );
	const path = require( "path" );
	const { execSync } = require( "child_process" );
	// const global_package_path = process.argv[ 0 ].split( "/bin/node" )[ 0 ] + "/lib/node_modules";
	const global_package_path = execSync( "npm root -g" ).toString().trim();
	const puppeteer = require( path.join( global_package_path ,  "puppeteer" ) );
	const fs = require( "fs" ).promises;
	function sleep( ms ) { return new Promise( resolve => setTimeout( resolve , ms ) ); }

	// Config
	const width = 3840; // 4K width
	const height = 2160; // 4K height
	const scale = 2; // 2x zoom ???
	const padding = 5;
	let input_file_path = process.argv[ 2 ];
	let output_file_type = "jpeg"
	let output_file_path = "output.jpeg"
	if ( process.argv.length > 3 ) { output_file_path = process.argv[ 3 ]; }

	// Runtime Config
	const interpolated_width = ( width + width / scale ) / 2
	const interpolated_height = ( height + height / scale ) / 2

	console.log( `Width === ${width}` );
	console.log( `Height === ${height}` );
	console.log( `Scale === ${scale}` );
	console.log( `Interpolated Width === ${interpolated_width}` );
	console.log( `Interpolated Height === ${interpolated_height}` );

	// Setup Puppeteer
	const browser = await puppeteer.launch({
		headless: "new" ,
		args: [
			"--disable-web-security" ,
			"--no-sandbox" ,
			"--disable-setuid-sandbox" ,
		] ,
		executablePath: "/usr/bin/chromium" ,
	});
	const page = await browser.newPage();
	let url;
	let svg_html_text;
	if ( input_file_path.startsWith( "http://" ) || input_file_path.startsWith( "https://" ) ) {
		url = input_file_path
		svg_html_text = `const response = await fetch("${input_file_path}");
		const svgText = await response.text();
		document.getElementById('svg-container').innerHTML = svgText;`;
		console.log( "using remote .svg file" , url );
	} else {
		url = `file://${input_file_path}`;
		const fileContent = await fs.readFile( input_file_path , "utf-8" );
		svg_html_text = `const svgText = \`${fileContent}\`; document.getElementById('svg-container').innerHTML = svgText;`;
		console.log( "using local .svg file" , url );
	}
	const html_content = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>SVG Test</title>
	<style>
	body, html {
	  margin: 0;
	  padding: 0;
	  width: 100%;
	  height: 100%;
	}
	</style>
</head>
<body>
	<br>
	<br>
	<div id="svg-container"></div>
	<script>
		( async ()=> {
			${svg_html_text} // Inject the svgText here
			document.getElementById( "svg-container" ).innerHTML = svgText;
			const svgRoot = document.querySelector( "svg" );
			const bbox = svgRoot.getBBox();
			const rect = svgRoot.getBoundingClientRect();
			svgRoot.setAttribute( "viewBox" , \`\${bbox.x} \${bbox.y} \${bbox.width} \${bbox.height}\` );
			svgRoot.setAttribute( "width" , "100%" );
			svgRoot.setAttribute( "height" , "100%" );
		})();
	</script>
</body>
</html>`;
	await page.setContent( html_content );

	// Give it some time to ensure that all fonts and images are loaded
	await page.waitForFunction( () => {
		let svg_root = document.querySelector( "svg" );
		return svg_root && svg_root.getBBox();
	});

	// This has to run in this order
	let bbox = await page.evaluate( ( padding ) => {
		let svg_root = document.querySelector( "svg" );
		let bbox = svg_root.getBBox();
		svg_root.setAttribute( "viewBox" , `${bbox.x - padding} ${bbox.y - padding} ${bbox.width + padding} ${bbox.height + padding}` );
		let x = document.querySelector( "svg" ).getBBox();
		return {
			x: x.x ,
			y: x.y ,
			width: x.width ,
			height: x.height
		};
	} , padding );
	console.log( bbox );

	await page.setViewport({
	  width: Math.floor( interpolated_width ) ,
	  height: Math.floor( interpolated_height ) ,
	  deviceScaleFactor: scale
	});

	const svg_element = await page.$( "svg" );
	await svg_element.screenshot( { path: output_file_path , type: output_file_type , quality: 100 } );
	await browser.close();
})();