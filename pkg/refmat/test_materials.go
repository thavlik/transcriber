package refmat

// TestReferenceMaterials is a list of reference materials that can be used for
// testing. This list is not exhaustive. It is meant to be used for testing
// purposes only. This is a good test video: https://youtu.be/gd4-FV_lwSE
var TestReferenceMaterials = []*ReferenceMaterial{{
	Terms: []string{
		"vertebral arch",
	},
	Images: []string{
		"https://refmat.nyc3.digitaloceanspaces.com/lumbar-vertebra-vertebral-arch-superior-view-745x550.png",
		"https://refmat.nyc3.digitaloceanspaces.com/General-Structure-of-a-Vertebrae.jpg",
	},
}, {
	Terms: []string{
		"ligamentum flavum",
		"ligamentum",
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
