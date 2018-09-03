var wasm_start, wasm_callback, go;

if (hvue_wasm) {
	if (!WebAssembly.instantiateStreaming) { // polyfill
		WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		};
	}

	let mod, inst;

	wasm_start = async function(file) {
		response = await fetch("/examples/wasm_exec.js")
		if(response.ok) {
			const source = await (await response).text();
		   eval(source);
		   go = new Go();
		} else {
			throw new Error('Network response was not ok.');
		}

		WebAssembly.instantiateStreaming(fetch("/examples/wasm/"+file+".wasm"), go.importObject).then((result) => {
			mod = result.module;
			inst = result.instance;
			go.run(inst);
		});
	}

	wasm_call_with_this = function(f) {
		return function() {
			f(this, ...arguments);
		}
	}

} else {
	// UNTESTED

	wasm_start = async function(file) {
		// OMG FIXME
		response = await fetch("/examples/"+file+"/"+file+".js")
		if(response.ok) {
			const source = await response.text();
			eval(source);
		} else {
			throw new Error('Network response was not ok.');
		}
	}

	wasm_call_with_this = function(f) {
		return f;
	}
}

function wasm_return_thing(thing) {
	return function() {
		return thing
	}
}

function wasm_new_data_func(templateObj, f) {
	return function() {
		newO = Object.assign({}, templateObj)
		f(newO) // runs later
		return newO
	}
}

