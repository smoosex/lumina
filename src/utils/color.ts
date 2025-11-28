/**
 * Parses an OKLCH color string and returns a hex color string.
 */
export function oklchToHex(oklchStr: string): string {
	// Captures: 1: L-value, 2: L-unit(%), 3: C-value, 4: H-value
	const match = oklchStr.match(
		/oklch\(\s*([\d.]+)(%?)\s+([\d.]+)\s+([\d.]+)(?:\s*\/\s*[\d.]+%?)?\s*\)/i,
	);

	if (!match) {
		// Fallback check
		if (oklchStr.startsWith("#") || oklchStr.startsWith("rgb"))
			return oklchStr;
		return "#000000";
	}

	let l = parseFloat(match[1] || "0");
	const isPercent = match[2] === "%";
	let c = parseFloat(match[3] || "0");
	let h = parseFloat(match[4] || "0");

	if (isPercent) l /= 100;

	// OKLCH -> OKLab
	const hRad = h * (Math.PI / 180);
	const a = c * Math.cos(hRad);
	const b = c * Math.sin(hRad);

	// OKLab -> Linear sRGB
	const l_ = l + 0.3963377774 * a + 0.2158037573 * b;
	const m_ = l - 0.1055613458 * a - 0.0638541728 * b;
	const s_ = l - 0.0894841775 * a - 1.291485548 * b;

	const l__ = l_ * l_ * l_;
	const m__ = m_ * m_ * m_;
	const s__ = s_ * s_ * s_;

	let rLin = 4.0767416621 * l__ - 3.3077115913 * m__ + 0.2309699292 * s__;
	let gLin = -1.2684380046 * l__ + 2.6097574011 * m__ - 0.3413193965 * s__;
	let bLin = -0.0041960863 * l__ - 0.7034186147 * m__ + 1.707614701 * s__;

	// Clipping before gamma is safer for simple hex conversion
	const clip = (x: number) => Math.max(0, Math.min(1, x));

	rLin = clip(rLin);
	gLin = clip(gLin);
	bLin = clip(bLin);

	// Linear sRGB -> Gamma sRGB
	const gamma = (x: number) => {
		return x >= 0.0031308
			? 1.055 * Math.pow(x, 1.0 / 2.4) - 0.055
			: 12.92 * x;
	};

	const r = Math.round(gamma(rLin) * 255);
	const g = Math.round(gamma(gLin) * 255);
	const blue = Math.round(gamma(bLin) * 255);

	// Convert to Hex
	const toHex = (n: number) => n.toString(16).padStart(2, "0");

	return `#${toHex(r)}${toHex(g)}${toHex(blue)}`;
}
