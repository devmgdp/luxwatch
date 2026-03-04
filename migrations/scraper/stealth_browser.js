const { chromium } = require("playwright-extra");
const stealth = require("puppeteer-extra-plugin-stealth")();
chromium.use(stealth);

async function getDeals(url) {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  try {
    await page.goto(url, { waitUntil: "networkidle", timeout: 60000 });

    const products = await page.evaluate(() => {
      const cards = document.querySelectorAll(".poly-card, .ui-search-result, .promotion-item, .items_container");
      
      return Array.from(cards).map((item) => {
        const name = item.querySelector("h2, .poly-component__title, .promotion-item__title")?.innerText;

        // Foca no container de preço ATUAL para evitar pegar o preço riscado
        const priceContainer = item.querySelector(".poly-price__current .andes-money-amount, .ui-search-price__second-line .andes-money-amount, .promotion-item__price .andes-money-amount");
        
        let price = null;
        if (priceContainer) {
          const fraction = priceContainer.querySelector(".andes-money-amount__fraction")?.innerText || "0";
          price = fraction.replace(/\./g, ""); // Remove pontos de milhar
        } else {
          // Fallback
          const fallback = item.querySelector(".andes-money-amount__fraction")?.innerText;
          price = fallback ? fallback.replace(/\./g, "") : null;
        }

        const imgElement = item.querySelector("img");
        const img = imgElement?.getAttribute("data-src") || imgElement?.src;
        const link = item.querySelector("a")?.href;

        return { name, price, image: img, link: link };
      }).filter((p) => p.name && p.price && p.link);
    });

    process.stdout.write(JSON.stringify(products.slice(0, 100)));
  } catch (err) {
    process.stderr.write(err.message);
  } finally {
    await browser.close();
  }
}
getDeals(process.argv[2]);
