package refmat

// TestReferenceMaterials is a list of reference materials that can be used for
// testing. This list is not exhaustive. It is meant to be used for testing
// purposes only. This is a good test video: https://youtu.be/gd4-FV_lwSE
var TestReferenceMaterials = []*ReferenceMaterial{{
	Terms: []string{
		"supraspinous ligament",
		"intertransverse ligament",
		"posterior longitudinal ligament",
		"anterior longitudinal ligament",
		"interspinous ligament",
		"facet capsular ligament",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/3746-ligaments_labeled.jpg",
		"https://refmat.nyc3.digitaloceanspaces.com/spu_article_asset_72ffb514b29c7d724ddbc64460a1b2cd933e9286.webp",
		"https://refmat.nyc3.digitaloceanspaces.com/Ligament-Injuries-1.jpg",
	},
}, {
	Terms: []string{"spinous process"},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/ok29q1av3b.jpg",
	},
}, {
	Terms: []string{
		"disc herniation",
		"herniated disc",
		"disc degeneration",
		"degenerated disc",
		"disc prolapse",
		"prolapsed disc",
		"disc extrusion",
		"extruded disc",
		"sequestered disc",
		"disc sequestration",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/spu_article_asset_59ce747548988425481ac3684e917c63ee81b142.webp",
		"https://refmat.nyc3.digitaloceanspaces.com/stages-of-disc-herniation.jpg",
	},
}, {
	Terms: []string{
		"pedicle",
		"pedicles",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/CervicalAnatomy-C3C4C5C6.webp",
	},
}, {
	Terms:  []string{"spinal canal"},
	Images: []string{"https://refmat.nyc3.digitaloceanspaces.com/aci3639_460x300.jpg"},
}, {
	Terms: []string{
		"vertebral arch",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/lumbar-vertebra-vertebral-arch-superior-view-745x550.png",
		"https://refmat.nyc3.digitaloceanspaces.com/General-Structure-of-a-Vertebrae.jpg",
	},
}, {
	Terms: []string{
		"ligamentum",
		"ligamenta",
		"ligamentum flavum",
		"ligamenta flavum",
		"flavum",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/ligamentum-flavum-1024x670.jpg",
		"https://refmat.nyc3.digitaloceanspaces.com/LigamentumFlavum.png",
	},
}, {
	Terms: []string{
		"facet joint",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/facet_joints_related_spine_structures_shutterstock_157672247.jpg",
		"https://refmat.nyc3.digitaloceanspaces.com/Thoracic-Facet-Syndrome.jpg",
	},
}}
