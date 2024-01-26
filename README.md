# e-Shadowing Transcriber
![Example Screenshot](app/captures/cap_00.jpg)

This is a prototype tool that transcribes an audio stream, detects the transcript's medical keywords, and automatically fetches supplementary content (definitions, pictures, drug data) when the user taps/clicks on a keyword.

It is meant to enhance the e-Shadowing experience by streamlining the process of finding good reference material without diverting one's attention away from the physician. When e-shadowing, I sometimes find myself searching Google for anatomy cartoons to refresh my memory, which is relatively cumbersome when compared to this tool.

This project serves more as an archive for an idea that I envisioned, rather than a sincere effort at producing a commercially viable product. It does exactly what I intended it to do, and I still strongly believe the concept has merit. I love medical education and seek to further our understanding of how technology may assist students and professors alike.

## Features
- Automatic recognition of key terms
    - Diagnoses
    - Anatomy
    - Organ systems
    - Medications (both generic and brand name)
- Context-sensitive image search results
    - Images are curated for relevance
    - Show only radiology images (anatomy)
    - Show only histology images (diagnoses)
- Automatic retrieval of drug data
    - Scrapes [DrugBank](https://go.drugbank.com/) with a [headless web browser](https://pptr.dev/)
    - Interactive 3D visualization of chemical structure
- Accurate transcription of medical language
    - Made possible with [Amazon Transcribe Medical](https://docs.aws.amazon.com/transcribe/latest/dg/transcribe-medical.html)
    - English only
- Simple and versatile architecture
    - Audio source is an [RTMP](https://en.wikipedia.org/wiki/Real-Time_Messaging_Protocol) stream
    - Developed to support [OBS](https://obsproject.com/)
    - Highly scalable [kubernetes](https://kubernetes.io/) backend

### Platform support
| Platform      | Supported          |
| ------------- | ------------------ |
| Windows 10/11 | :x: |
| MacOS         | :heavy_check_mark: |
| Linux         | :heavy_check_mark: |
| iOS           | :heavy_check_mark: |
| Android       | :heavy_check_mark: |
| Chrome        | :heavy_check_mark: |

## Design
The application is packaged as a [Flutter](https://flutter.dev/) front-end and a [helm](https://helm.sh/) chart backend for kubernetes. The application should scale horizontally to accomodate hundreds (or even thousands) of concurrent users, without any further optimization.

## Cost Warning
Running this tool can be financially expensive. If you leave it running for an hour, don't be surprised if your AWS bill is in the hundreds of dollars. For this reason, I'm apprehensive to believe such a tool has any commercial viability at all.

## License
All code in this repository is released under [MIT](LICENSE-MIT) / [Apache 2.0](LICENSE-Apache) dual license, which is extremely permissive. Please open an issue if somehow these terms are insufficient.