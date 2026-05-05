(() => {
	const DAY = 'day';
	const MONTH = 'month';
	const YEAR = 'year';
	const HOLD_DELAY_MS = 400;
	const HOLD_REPEAT_MS = 120;

	const isValidDateObject = (value) => value instanceof Date && !Number.isNaN(value.getTime());
	const cloneDate = (date) => new Date(date.getFullYear(), date.getMonth(), date.getDate());
	const daysInMonth = (year, monthIndex) => new Date(year, monthIndex + 1, 0).getDate();

	const formatYyyymmdd = (date) => {
		const y = String(date.getFullYear()).padStart(4, '0');
		const m = String(date.getMonth() + 1).padStart(2, '0');
		const d = String(date.getDate()).padStart(2, '0');
		return `${y}${m}${d}`;
	};

	const parseYyyymmdd = (value) => {
		if (typeof value !== 'string' || !/^\d{8}$/.test(value)) return null;
		const year = Number.parseInt(value.slice(0, 4), 10);
		const month = Number.parseInt(value.slice(4, 6), 10);
		const day = Number.parseInt(value.slice(6, 8), 10);
		if (month < 1 || month > 12) return null;
		const maxDay = daysInMonth(year, month - 1);
		if (day < 1 || day > maxDay) return null;
		const date = new Date(year, month - 1, day);
		if (!isValidDateObject(date)) return null;
		if (date.getFullYear() !== year || date.getMonth() !== month - 1 || date.getDate() !== day) return null;
		return date;
	};

	const normalizeDateValue = (value) => {
		if (isValidDateObject(value)) return cloneDate(value);
		if (typeof value === 'string') return parseYyyymmdd(value);
		return null;
	};

	const resolveDateOption = (value, fallback) => {
		if (value == null || value === '') return cloneDate(fallback);
		const parsed = normalizeDateValue(value);
		if (!parsed) return cloneDate(fallback);
		return parsed;
	};

	const clampDate = (date, minDate, maxDate) => {
		if (date < minDate) return cloneDate(minDate);
		if (date > maxDate) return cloneDate(maxDate);
		return date;
	};

	class QuoteStepDateControl {
		constructor(targetElement, options = {}) {
			if (!(targetElement instanceof HTMLElement)) return;
			const today = new Date();
			const defaultMinDate = new Date(today.getFullYear() - 100, today.getMonth(), today.getDate());
			const defaultMaxDate = new Date(today.getFullYear() + 100, today.getMonth(), today.getDate());

			this.minDate = resolveDateOption(options.dateMin, defaultMinDate);
			this.maxDate = resolveDateOption(options.dateMax, defaultMaxDate);
			if (this.minDate > this.maxDate) this.maxDate = cloneDate(this.minDate);

			const initial = normalizeDateValue(options.value) || cloneDate(today);
			this.currentDate = clampDate(initial, this.minDate, this.maxDate);
			this.onChange = typeof options.onChange === 'function' ? options.onChange : null;

			this.holdDelayTimer = null;
			this.holdRepeatTimer = null;
			this.activeHoldMeta = null;
			this.suppressNextClick = false;

			this.host = targetElement;
			this.host.classList.add('date-control-host');
			this.root = this.host;
			this.render();
			this.syncInputs();
		}

		getValue() { return formatYyyymmdd(this.currentDate); }
		getDate() { return cloneDate(this.currentDate); }

		setValue(value) {
			const normalized = normalizeDateValue(value);
			if (!normalized) return false;
			const next = clampDate(cloneDate(normalized), this.minDate, this.maxDate);
			const changed = next.getTime() !== this.currentDate.getTime();
			this.currentDate = next;
			this.syncInputs();
			if (changed && this.onChange) this.onChange(this.getValue(), this.getDate());
			return true;
		}

		render() {
			const wrapper = document.createElement('div');
			wrapper.className = 'date-control';
			this.dayInput = this.buildSegment(wrapper, DAY, 2);
			this.monthInput = this.buildSegment(wrapper, MONTH, 2);
			this.yearInput = this.buildSegment(wrapper, YEAR, 4);
			this.root.replaceChildren(wrapper);
		}

		buildSegment(wrapper, segmentType, maxLen) {
			const segment = document.createElement('div');
			segment.className = `segment ${segmentType}`;

			const upBtn = document.createElement('button');
			upBtn.className = 'arrow-btn up';
			upBtn.type = 'button';
			upBtn.setAttribute('aria-label', `Increase ${segmentType}`);
			upBtn.addEventListener('click', () => {
				if (this.suppressNextClick) { this.suppressNextClick = false; return }
				if (this.activeHoldMeta) return;
				this.stepSegment(segmentType, 1);
			});
			this.attachHoldBehavior(upBtn, segmentType, 1);

			const input = document.createElement('input');
			input.className = 'value-input';
			input.type = 'text';
			input.inputMode = 'numeric';
			input.maxLength = maxLen;
			input.readOnly = true;
			input.autocomplete = 'off';
			input.spellcheck = false;
			input.setAttribute('aria-label', segmentType);
			input.addEventListener('keydown', (event) => {
				if (event.key === 'Tab') return;
				event.preventDefault();
			});

			const downBtn = document.createElement('button');
			downBtn.className = 'arrow-btn down';
			downBtn.type = 'button';
			downBtn.setAttribute('aria-label', `Decrease ${segmentType}`);
			downBtn.addEventListener('click', () => {
				if (this.suppressNextClick) { this.suppressNextClick = false; return }
				if (this.activeHoldMeta) return;
				this.stepSegment(segmentType, -1);
			});
			this.attachHoldBehavior(downBtn, segmentType, -1);

			segment.append(upBtn, input, downBtn);
			wrapper.appendChild(segment);
			return input;
		}

		syncInputs() {
			this.dayInput.value = `${String(this.currentDate.getDate()).padStart(2, '0')}.`;
			this.monthInput.value = `${String(this.currentDate.getMonth() + 1).padStart(2, '0')}.`;
			this.yearInput.value = String(this.currentDate.getFullYear()).padStart(4, '0');
			this.host.setAttribute('value', this.getValue());
			this.host.setAttribute('data-min', formatYyyymmdd(this.minDate));
			this.host.setAttribute('data-max', formatYyyymmdd(this.maxDate));
		}

		clearHoldTimers() {
			if (this.holdDelayTimer) { clearTimeout(this.holdDelayTimer); this.holdDelayTimer = null }
			if (this.holdRepeatTimer) { clearInterval(this.holdRepeatTimer); this.holdRepeatTimer = null }
		}

		stepWithHoldAcceleration(segmentType, direction, accelerated) {
			if (segmentType === YEAR && accelerated) { this.stepSegment(YEAR, direction * 5); return }
			this.stepSegment(segmentType, direction);
		}

		startHold(segmentType, direction, pointerId) {
			this.clearHoldTimers();
			this.activeHoldMeta = { segmentType, direction, pointerId };
			this.suppressNextClick = false;
			this.holdDelayTimer = setTimeout(() => {
				if (!this.activeHoldMeta) return;
				this.suppressNextClick = true;
				this.stepWithHoldAcceleration(segmentType, direction, true);
				this.holdRepeatTimer = setInterval(() => {
					if (!this.activeHoldMeta) return;
					this.stepWithHoldAcceleration(segmentType, direction, true);
				}, HOLD_REPEAT_MS);
			}, HOLD_DELAY_MS);
		}

		endHold(pointerId) {
			if (!this.activeHoldMeta) return;
			if (pointerId != null && this.activeHoldMeta.pointerId !== pointerId) return;
			this.clearHoldTimers();
			this.activeHoldMeta = null;
		}

		attachHoldBehavior(button, segmentType, direction) {
			button.addEventListener('pointerdown', (event) => {
				if (event.button !== 0) return;
				event.preventDefault();
				button.setPointerCapture(event.pointerId);
				this.startHold(segmentType, direction, event.pointerId);
			});
			button.addEventListener('pointerup', (event) => { this.endHold(event.pointerId); });
			button.addEventListener('pointercancel', (event) => { this.endHold(event.pointerId); });
			button.addEventListener('lostpointercapture', (event) => { this.endHold(event.pointerId); });
		}

		stepSegment(segmentType, delta) {
			const y = this.currentDate.getFullYear();
			const m = this.currentDate.getMonth();
			const d = this.currentDate.getDate();
			let nextDate;
			if (segmentType === DAY) {
				nextDate = new Date(y, m, d + delta);
			} else if (segmentType === MONTH) {
				const nextMonth = m + delta;
				const nextYear = y + Math.floor(nextMonth / 12);
				const normalizedMonth = ((nextMonth % 12) + 12) % 12;
				const safeDay = Math.min(d, daysInMonth(nextYear, normalizedMonth));
				nextDate = new Date(nextYear, normalizedMonth, safeDay);
			} else {
				const nextYear = y + delta;
				const safeDay = Math.min(d, daysInMonth(nextYear, m));
				nextDate = new Date(nextYear, m, safeDay);
			}
			this.setValue(nextDate);
		}
	}

	const mountDateControl = (root) => {
		if (!(root instanceof HTMLElement)) return;
		if (root.getAttribute('data-date-init') === '1') return;
		const hidden = root.querySelector('input[data-date-hidden="1"][name]');
		const host = root.querySelector('[data-date-host="1"]');
		if (!(hidden instanceof HTMLInputElement) || !(host instanceof HTMLElement)) return;

		const currentValue = hidden.value || host.getAttribute('value') || '';
		const dateMin = hidden.getAttribute('data-min') || host.getAttribute('data-min') || '';
		const dateMax = hidden.getAttribute('data-max') || host.getAttribute('data-max') || '';

		new QuoteStepDateControl(host, {
			value: currentValue,
			dateMin,
			dateMax,
			onChange: (value) => {
				if (hidden.value === value) return;
				hidden.value = value;
				hidden.dispatchEvent(new Event('input', { bubbles: true }));
				hidden.dispatchEvent(new Event('change', { bubbles: true }));
			},
		});
		root.setAttribute('data-date-init', '1');
	};

	const initDateControls = (scope = document) => {
		const base = scope instanceof Document || scope instanceof HTMLElement ? scope : document;
		for (const root of base.querySelectorAll('[data-date-control="1"]')) mountDateControl(root);
	};

	window.QuoteDateControlInit = initDateControls;
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', () => { initDateControls(document); });
	} else {
		initDateControls(document);
	}
})();
