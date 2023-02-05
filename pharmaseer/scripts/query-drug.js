// Run with:
// docker run -i --init --cap-add=SYS_ADMIN --rm ghcr.io/puppeteer/puppeteer:latest node -e "$(cat query-drug.js)"
const puppeteer = require('puppeteer');
//const {TimeoutError} = require('puppeteer/Errors');
//try {
//  await page.waitForSelector(yourSelector, {timeout: 5000});
//} catch (e) {
//  if (e instanceof TimeoutError) {
//      // Do something if this is a timeout.
//  }
//}

const TIMEOUT = 5000;

function parseWeight(input) {
    // "Average: 321.158 Monoisotopic: 320.011932988 "
    let avgIndex = input.indexOf('Average: ');
    let end = input.indexOf(' ', avgIndex + 9);
    if (end == -1) {
        end = input.length;
    }
    let avgValue = input.substring(avgIndex + 9, end).trim();
    let monoIndex = input.indexOf('Monoisotopic: ');
    end = input.indexOf(' ', monoIndex + 14);
    if (end == -1) {
        end = input.length;
    }
    let monoValue = input.substring(monoIndex + 14, end).trim();
    return {
        average: avgValue,
        monoisotopic: monoValue,
    };
}

const getTextValue = async (page, id) => {
    const el = await page.waitForSelector(id, { timeout: TIMEOUT });
    const val = await page.evaluateHandle(el => el.nextElementSibling, el);
    const text = await (await val.getProperty('textContent')).jsonValue();
    return text.trim();
}

const getArrayElements = async (page, selector) => {
    const el = await page.waitForSelector(selector, { timeout: TIMEOUT });
    const items = (await (await page.evaluateHandle(el => {
        const next = el.nextElementSibling;
        const data = [];
        for (const child of next.children) {
            data.push(child.innerText);
        }
        return data;
    }, el)).jsonValue())[0].split('\n');
    return items;
};

const getOverview = async (page) => {
    const summary = await getTextValue(page, '#summary');
    const brandNames = (await getTextValue(page, '#brand-names')).split(', ');
    const genericName = await getTextValue(page, '#generic-name');
    const drugbankAccessionNumber = await getTextValue(page, '#drugbank-accession-number');
    const background = await getTextValue(page, '#background');
    const groups = (await getTextValue(page, '#groups')).split(', ');
    const structureImageEl = await page.waitForSelector('.structure a img', { timeout: TIMEOUT });
    const structureImage = await (await page.evaluateHandle(el => el.getAttribute('src'), structureImageEl)).jsonValue();

    const externalIDs = await getArrayElements(page, '#external-ids');

    await page.click('#structure-download');
    const [pdbEl] = await page.$x("//a[contains(text(), 'PDB')]");
    const pdb = await (await pdbEl.getProperty('href')).jsonValue();

    const weight = parseWeight(await getTextValue(page, '#weight'));
    const chemicalFormula = await getTextValue(page, '#chemical-formula');
    const synonymns = await getArrayElements(page, '#synonyms');

    return {
        summary,
        brandNames,
        genericName,
        drugbankAccessionNumber,
        background,
        groups,
        weight,
        chemicalFormula,
        synonymns,
        structure: {
            image: 'https://go.drugbank.com' + structureImage,
            pdb,
        },
        externalIDs,
    };
}

const getMetabolism = async (page) => {
    const el = await page.waitForSelector('#metabolism', { timeout: TIMEOUT });
    const val = await page.evaluateHandle(el => el.nextElementSibling, el);
    let description = await (await val.getProperty('textContent')).jsonValue();
    let i = description.indexOf('Hover over products below to view reaction partners');
    if (i != -1) {
        description = description.substring(0, i);
    }
    //const tree = await page.evaluateHandle(el => {
    //    // ul.metabolite-pathway
    //    const parseListItem = (el) => {
    //        el.querySelectorAll('> span').forEach((span) => {});
    //        el.querySelectorAll('> li').forEach((li) => {});
    //    };
    //    const table = el.querySelector('ul.metabolite-pathway');
    //    return parseListItem(table);
    //}, val);
    return {
        description,
    };
};

const getReferences = async (page) => {
    const el = await page.waitForSelector('#general-references', { timeout: TIMEOUT });
    const val = await page.evaluateHandle(el => el.nextElementSibling, el);
    const general = await (await page.evaluateHandle(el => {
        const tds = Array.from(el.querySelectorAll('ol.cite-this-references > li'))
        return tds.map((td, i) => {
            let title = td.innerText;
            title = title.replaceAll('[Article]', '');
            title = title.replaceAll('[Link]', '');
            title = title.trim();
            const out = {
                index: i + 1,
                title,
            }
            const a = td.querySelector('a');
            if (a != null) {
                let link = a.getAttribute('href');
                if (link.startsWith('/')) {
                    link = 'https://go.drugbank.com' + link;
                }
                out['link'] = link;
            }
            return out;
        });
    }, val)).jsonValue();
    return {
        general,
    };
};

const getPharmacology = async (page) => {
    const indication = await getTextValue(page, '#indication');

    const associatedConditionsEl = await page.waitForSelector('#associated-conditions', { timeout: TIMEOUT });
    const associatedConditions = (await (await page.evaluateHandle(el => {
        const next = el.nextElementSibling;
        const data = [];
        for (const child of next.children) {
            data.push(child.innerText);
        }
        return data;
    }, associatedConditionsEl)).jsonValue())[0].split('\n');


    //const associatedTherapiesEl = await page.waitForSelector('#associated-therapies');
    //const associatedTherapies = (await (await page.evaluateHandle(el => {
    //    const next = el.nextElementSibling;
    //    const data = [];
    //    for (const child of next.children) {
    //        data.push(child.innerText);
    //    }
    //    return data;
    //}, associatedTherapiesEl)).jsonValue())[0].split('\n');

    const pharmacodynamics = await getTextValue(page, '#pharmacodynamics');
    const mechanismOfAction = await getTextValue(page, '#mechanism-of-action');
    const absorption = await getTextValue(page, '#absorption');
    const volumeOfDistribution = await getTextValue(page, '#volume-of-distribution');
    const proteinBinding = await getTextValue(page, '#protein-binding');
    const metabolism = await getMetabolism(page);
    const routeOfElimination = await getTextValue(page, '#route-of-elimination');
    const halfLife = await getTextValue(page, '#half-life');
    const clearance = await getTextValue(page, '#clearance');
    const toxicity = await getTextValue(page, '#toxicity');

    return {
        indication,
        associatedConditions,
        //associatedTherapies,
        pharmacodynamics,
        mechanismOfAction,
        absorption,
        volumeOfDistribution,
        proteinBinding,
        metabolism,
        routeOfElimination,
        halfLife,
        clearance,
        toxicity,
    };
}

const parsePropertyTable = async (page, selector) => {
    const el = await page.waitForSelector(selector, { timeout: TIMEOUT });
    const val = await page.evaluateHandle(el => el.nextElementSibling, el);
};

const getProperties = async (page) => {
    const experimentalProperties = await parsePropertyTable(page, '#experimental-properties');
    const predictedProperties = await parsePropertyTable(page, '#predicted-properties');
    const predictedADMETFeatures = await parsePropertyTable(page, '#predicted-admet-features');
    return {
        experimentalProperties,
        predictedProperties,
        predictedADMETFeatures,
    }
};

(async () => {
    const url = process.env.INPUT_URL;
    if (!url) {
        throw new Error("missing INPUT_URL");
    }

    const browser = await puppeteer.launch();
    try {
        const page = await browser.newPage();
        await page.goto(url);

        const general = await getOverview(page);

        // get the pharmacology data
        //await page.click('#pharmacology-sidebar-header');
        //const pharmacology = await getPharmacology(page);

        //await page.click('#references-sidebar-header');
        //const references = await getReferences(page);

        //await page.click('#properties-sidebar-header');
        //const properties = await parseProperties(page);

        console.log(JSON.stringify({
            ...general,
            //pharmacology,
            //references,
        }, null, 2));
    } finally {
        await browser.close();
    }
})();

