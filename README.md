# ASCII Art Web

## Description
Ascii-art-web consists in creating and running a server, in which it will be possible to use a web GUI (graphical user interface) version of ascii-art project.
ASCII Art Web is a web-based application that converts text input into visually appealing ASCII art. Users can input their desired text, select a bannner, color, text align, and view the corresponding ASCII art output directly on the website. This project demonstrates the seamless integration of backend algorithms and frontend interfaces to provide a creative text-to-ASCII experience.

---

## Authors
- **[Parisa Rahimi Darabad]** - Senior Backend Developer  
- **[Majid Rouhani]** - Senior Backend Developer  

---

## Usage
### How to Run
1. Clone this repository to your local machine:
   ```bash
   git clone https://gitea.com/mrouhani/ascii-art-web.git

2. Navigate to the project directory:
    ```bash
    cd ascii-art-web

3. Start the development server:
    ```bash
    go run .

4. Open your browser and navigate to:
    ```arduino
    http://localhost:8080

5. Input your text, choose a banner, color, align, and generate your ASCII art!

4. You can run test from root with this command:
    ```arduino
    go test ./...


## Project Structure and Implementation
Project has 2 main components

Backend: Include webserver, Api, Tests, AsciiArt service and Banners

Frontend: Include html templates, error files and assets

In the root of the project there is a runner file to run project from the root while keep structure of the project organized.