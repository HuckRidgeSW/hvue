var wasm_start;
var go;

console.log("hvue_wasm:", hvue_wasm)

if (hvue_wasm) {
	if (!WebAssembly.instantiateStreaming) { // polyfill
		WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		};
	}

	let mod, inst;

	wasm_start = async function(name) {
		response = await fetch("/vendor/wasm_exec.js")
		if(response.ok) {
			const source = await (await response).text();
		   eval(source);
		   go = new Go();
		} else {
			throw new Error('Network response was not ok.');
		}

		WebAssembly.instantiateStreaming(fetch("/examples/"+name+"/"+name+".wasm"), go.importObject).then((result) => {
			mod = result.module;
			inst = result.instance;
			go.run(inst);
		});
	}

} else {
	wasm_start = async function(name) {
		response = await fetch("/examples/"+name+"/"+name+".js")
		if(response.ok) {
			const source = await response.text();
			eval(source);
		} else {
			throw new Error('Network response was not ok.');
		}
	}
}

