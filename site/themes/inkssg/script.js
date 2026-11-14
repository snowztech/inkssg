// The version next to the logo, from the latest GitHub release. Stays hidden
// if the request fails.
(async () => {
  const el = document.querySelector("[data-version]");
  if (!el) return;
  const key = "inkssg-version";
  let v;
  try { v = sessionStorage.getItem(key); } catch {}
  if (!v) {
    try {
      const res = await fetch("https://api.github.com/repos/snowztech/inkssg/releases/latest");
      if (!res.ok) return;
      v = (await res.json()).tag_name;
      try { sessionStorage.setItem(key, v); } catch {}
    } catch { return; }
  }
  el.textContent = v;
  el.hidden = false;
})();

// The hero demo: the heading in content.md is retyped, and the page on the
// right follows, the way inkssg serve reloads it. Still when motion is off.
(() => {
  const demo = document.querySelector("[data-demo]");
  if (!demo || matchMedia("(prefers-reduced-motion: reduce)").matches) return;
  const src = demo.querySelector("[data-src-title]");
  const out = demo.querySelector("[data-out-title]");
  const flash = demo.querySelector("[data-reload]");
  const titles = ["Hi, I'm Ada", "Small tools, built slowly", "Notes from the terminal"];
  const wait = (ms) => new Promise((r) => setTimeout(r, ms));
  let i = 0;
  // Only run while the demo is on screen and the tab is visible.
  let seen = false;
  new IntersectionObserver(([e]) => { seen = e.isIntersecting; }).observe(demo);
  const ready = async () => { while (!seen || document.hidden) await wait(400); };

  const loop = async () => {
    await wait(2600);
    await ready();
    const next = titles[++i % titles.length];
    while (src.textContent.length) { src.textContent = src.textContent.slice(0, -1); await wait(28); }
    for (const ch of next) { src.textContent += ch; await wait(55); }
    await wait(350);
    out.textContent = next;
    flash.classList.remove("on"); void flash.offsetWidth; flash.classList.add("on");
    loop();
  };
  demo.classList.add("is-live");
  loop();
})();

// Docs: a link on every section heading, and a copy button on every code block.
(() => {
  const prose = document.querySelector(".docs article.sn-prose");
  if (!prose) return;
  prose.querySelectorAll("h2[id]").forEach((h) => {
    const a = document.createElement("a");
    a.className = "anchor"; a.href = "#" + h.id; a.textContent = "#";
    a.setAttribute("aria-label", "Link to " + h.textContent);
    h.prepend(a);
  });
  prose.querySelectorAll("pre").forEach((pre) => {
    const b = document.createElement("button");
    b.className = "sn-copy-btn"; b.type = "button"; b.setAttribute("aria-label", "copy");
    b.addEventListener("click", async () => {
      try { await navigator.clipboard.writeText(pre.querySelector("code")?.textContent ?? pre.textContent); } catch { return; }
      b.dataset.copied = "true"; setTimeout(() => delete b.dataset.copied, 1200);
    });
    pre.append(b);
  });
})();
