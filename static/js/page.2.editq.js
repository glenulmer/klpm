(() => {
	const selector = '#EditQForm input[name], #EditQForm textarea[name], #EditQForm select[name], #EditQForm button[name]';
	const debounceMs = 240;
	const foldIds = ['EditQPreexCard', 'EditQDependantsCard', 'EditQReviewCard'];
	const foldDefaults = new Map([
		['EditQPreexCard', false],
		['EditQDependantsCard', false],
		['EditQReviewCard', true],
	]);
	const foldState = new Map();
	let seq = 0;
	const timers = new Map();
	let postChain = Promise.resolve();

	const captureFoldStates = () => {
		for (const id of foldIds) {
			const el = document.getElementById(id);
			if (!(el instanceof HTMLDetailsElement)) continue;
			foldState.set(id, el.open);
		}
	};

	const applyFoldStates = () => {
		for (const id of foldIds) {
			const el = document.getElementById(id);
			if (!(el instanceof HTMLDetailsElement)) continue;
			if (foldState.has(id)) {
				el.open = !!foldState.get(id);
				continue;
			}
			if (foldDefaults.has(id)) {
				el.open = !!foldDefaults.get(id);
			}
		}
	};

	const captureFocus = () => {
		const active = document.activeElement;
		if (!(active instanceof HTMLElement)) return null;
		if (!(active instanceof HTMLInputElement) && !(active instanceof HTMLTextAreaElement)) return null;
		const name = active.getAttribute('name') || '';
		if (!name) return null;
		const out = { name, start: 0, end: 0 };
		if (typeof active.selectionStart === 'number') out.start = active.selectionStart;
		if (typeof active.selectionEnd === 'number') out.end = active.selectionEnd;
		return out;
	};

	const restoreFocus = (focus) => {
		if (!focus) return;
		const list = document.querySelectorAll('#EditQForm [name]');
		let target = null;
		for (const el of list) {
			if (!(el instanceof HTMLInputElement) && !(el instanceof HTMLTextAreaElement)) continue;
			if ((el.getAttribute('name') || '') !== focus.name) continue;
			target = el;
			break;
		}
		if (!target) return;
		target.focus();
		if (typeof target.setSelectionRange === 'function') {
			target.setSelectionRange(focus.start, focus.end);
		}
	};

	const autosize = (el) => {
		if (!(el instanceof HTMLTextAreaElement)) return;
		if (!el.classList.contains('editq-grow-input')) return;
		el.style.height = 'auto';
		el.style.height = `${Math.max(el.scrollHeight, 34)}px`;
	};

	const autosizeAll = () => {
		const list = document.querySelectorAll('#EditQForm textarea.editq-grow-input');
		for (const el of list) autosize(el);
	};
	const initDateControls = (root = document) => {
		if (typeof window.QuoteDateControlInit !== 'function') return;
		window.QuoteDateControlInit(root);
	};

	const controlValue = (el) => {
		if (el instanceof HTMLButtonElement) return el.value || '1';
		if (el instanceof HTMLInputElement && el.type === 'checkbox') return el.checked ? '1' : '0';
		return el.value;
	};

	const postChange = (name, value) => {
		const call = ++seq;
		postChain = postChain
			.catch(() => {})
			.then(() => {
				const focus = captureFocus();
				captureFoldStates();
				const form = new FormData();
				form.append('name', name);
				form.append('value', value);

				return fetch('/quote-review-change', {
					method: 'POST',
					body: form,
					credentials: 'same-origin',
				})
					.then((res) => (res.ok ? res.json() : []))
					.then((messages) => {
						if (call !== seq || !Array.isArray(messages)) return;
						for (const msg of messages) {
							if (!msg || msg.kind !== 'rewrite') continue;
							const target = document.querySelector(msg.target);
							if (!target) continue;
							if (msg.method === 'outerHTML') target.outerHTML = msg.content;
							if (msg.method === 'innerHTML') target.innerHTML = msg.content;
						if (msg.method === 'remove') target.remove();
					}
					applyFoldStates();
					autosizeAll();
					initDateControls();
					restoreFocus(focus);
				});
			})
			.catch(() => {});
	};

	const sendChange = (name, value, force = false) => {
		if (force) { postChange(name, value); return }
		postChange(name, value);
	};

	const schedule = (name, value) => {
		const prev = timers.get(name);
		if (prev) window.clearTimeout(prev);
		const next = window.setTimeout(() => {
			timers.delete(name);
			sendChange(name, value);
		}, debounceMs);
		timers.set(name, next);
	};

	const onChange = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLElement)) return;
		if (!el.matches(selector)) return;
		if (el instanceof HTMLButtonElement) return;
		if (ev.type === 'input' && !(el instanceof HTMLTextAreaElement) && !(el instanceof HTMLInputElement && el.type === 'text')) return;
		const name = el.getAttribute('name') || '';
		if (!name) return;
		const value = controlValue(el);
		if (el instanceof HTMLTextAreaElement) {
			autosize(el);
			schedule(name, value);
			return;
		}
		if (el instanceof HTMLInputElement && el.type === 'text') {
			schedule(name, value);
			return;
		}
		sendChange(name, value);
	};

	const onClick = (ev) => {
		const target = ev.target;
		if (!(target instanceof HTMLElement)) return;
		const el = target.closest('#EditQForm button[name]');
		if (!(el instanceof HTMLButtonElement)) return;
		if (!el.matches(selector)) return;
		const name = el.getAttribute('name') || '';
		if (!name) return;
		if (name === 'DownloadExcel') {
			ev.preventDefault();
			const lang = document.querySelector('#EditQForm select[name="lang"]');
			const slim = document.querySelector('#EditQForm select[name="slim"]');
			const save = (key, value) => {
				const form = new FormData();
				form.append('name', key);
				form.append('value', value);
				return fetch('/quote-review-change', {
					method: 'POST',
					body: form,
					credentials: 'same-origin',
				});
			};
			const langVal = lang instanceof HTMLSelectElement ? lang.value : '';
			const slimVal = slim instanceof HTMLSelectElement ? slim.value : '0';
			Promise.all([
				save('lang', langVal),
				save('slim', slimVal),
			]).finally(() => {
				window.location.assign('/download-excel');
			});
			return;
		}
		ev.preventDefault();
		sendChange(name, controlValue(el), true);
	};

	const onFoldToggle = (ev) => {
		const el = ev.target;
		if (!(el instanceof HTMLDetailsElement)) return;
		if (!foldIds.includes(el.id)) return;
		foldState.set(el.id, el.open);
	};

	document.addEventListener('input', onChange);
	document.addEventListener('change', onChange);
	document.addEventListener('click', onClick);
	document.addEventListener('toggle', onFoldToggle, true);
	applyFoldStates();
	autosizeAll();
	initDateControls();
})();
