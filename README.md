# Product Tracker

Product Tracker is a Go-based application that allows you to track products in a database. It supports inserting new products and fetching products by name.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API](#api)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/product-tracker.git
    cd product-tracker
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Set up the database:

    Ensure you have PostgreSQL installed and running. Create a database and a table with the following structure:

    ```sql
    CREATE TABLE product_tracker (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        quantity INT NOT NULL,
        energy_consumed FLOAT NOT NULL,
        date DATE NOT NULL
    );
    ```

4. Configure the database connection:

    Update the `db` package to include your database connection details.

## Usage

1. Run the application:

    ```sh
    go run main.go
    ```

2. Use the provided API endpoints to interact with the application.

## API

### Insert Product

- **Endpoint**: `/insert`
- **Method**: `POST`
- **Request Body**:

    ```json
    {
        "name": "Product Name",
        "quantity": 10,
        "energy_consumed": 100.5,
        "date": "2025-03-10"
    }
    ```

- **Response**:

    ```json
    {
        "message": "Product inserted successfully"
    }
    ```

### Get Products by Name

- **Endpoint**: `/products`
- **Method**: `GET`
- **Query Parameters**:

    - `name`: The name of the product to fetch.

- **Response**:

    ```json
    [
        {
            "name": "Product Name",
            "quantity": 10,
            "energy_consumed": 100.5,
            "date": "2025-03-10"
        },
        ...
    ]
    ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
