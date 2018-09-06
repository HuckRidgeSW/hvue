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

function wasm_call_with_this(f) {
	return function() {
		f(this, ...arguments);
	}
}

function wasm_new_data_func(templateObj, f) {
	return function() {
		var newO;

		// Create a new object, based on the template
		newO = Object.assign({}, templateObj);
		// newO.hvue_vm = this;

		// Not sure I need this?
		if (newO.hvue_dataID === undefined) {
			var dataID = this.$parent.$data.hvue_dataID;
			if (dataID !== undefined) {
				newO.hvue_dataID = dataID
			}
		}

		// Call the hvue function to initialize these fields
		f(this, newO); // wasm: runs later; GopherJS: runs now

		return newO;
	}
}
