# e-Shadowing Transcriber
![Example Screenshot](app/captures/cap_00.jpg)

This is a tool that transcribes an audio stream, detects the transcript's medical keywords, and automatically fetches supplementary content (definitions, pictures, drug data, etc.) when the user taps/clicks on a keyword.

It is meant to enhance the e-Shadowing experience by streamlining the process of finding good reference material without diverting one's attention from what the physician is doing. When e-shadowing, I sometimes find myself searching Google for anatomy cartoons to refresh my memory, which is relatively cumbersome when compared to this tool.

This project serves more as an archive for an idea that I envisioned, rather than a sincere effort at producing a commercially viable product. This tool does exactly what I designed it to do, and I still strongly believe it has merit as a concept.

## Design
The application is packaged as a [Flutter](https://flutter.dev/) front-end and a [helm](https://helm.sh/) chart backend for kubernetes. The application should scale horizontally to accomodate hundreds (or even thousands) of concurrent users, without any further optimization.

## Cost Warning
Running this tool can be financially expensive. If you leave it running for an hour, don't be surprised if your AWS bill is in the hundreds of dollars. For this reason, I'm apprehensive to believe such a tool has any commercial viability at all.

## License
All code in this repository is released under [MIT](LICENSE-MIT) / [Apache 2.0](LICENSE-Apache) dual license, which is extremely permissive. Please open an issue if somehow these terms are insufficient.