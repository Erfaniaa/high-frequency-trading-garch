# Crypto High-Frequency Trading Using GARCH
Retrieve data from Binance and simulate high-frequency trading on them using the GARCH model

## Description
This project is a Go application designed to simulate a cryptocurrency trading strategy based on real-time market data and advanced statistical modeling. The primary goal is to implement and backtest a trading algorithm that uses the GARCH (Generalized Autoregressive Conditional Heteroskedasticity) model to predict market volatility and make informed trading decisions.

## Key Features
**Real-Time Data Retrieval from Binance:**

- WebSocket Connection: Establishes a live WebSocket connection to Binance's streaming API to receive real-time trade data for a specified trading pair (e.g., BTC/USDT).
- Live Price Updates: Continuously fetches the latest price information to ensure the simulation uses up-to-date market conditions.

**Volatility Calculation with GARCH Model:**
- GARCH(1,1) Implementation: Utilizes the GARCH(1,1) model to estimate future price volatility based on historical price data.
Sliding Window Analysis: Processes a moving window of recent price returns to input into the model for ongoing volatility prediction.
- Model Parameters: Configurable parameters (Omega, Alpha, Beta) allow for fine-tuning the model to better fit market behavior.

**Trading Strategy Simulation:**
- Limit Order Placement: Simulates placing limit buy and sell orders based on volatility predictions and calculated price levels.
Order Management: Handles order expiration, execution, and tracking over the course of the simulation.
- Wallet Management: Maintains a virtual wallet with USDT and cryptocurrency balances, updating holdings based on executed trades.
- Trade Sizing: Determines trade amounts as a percentage of the wallet, enabling risk management and portfolio scaling.

**Backtesting the Algorithm:**
- Performance Evaluation: Runs the trading strategy over a set number of iterations to assess profitability and effectiveness.
- Metrics Calculation: Computes key performance indicators such as mean spread percentage, filled orders percentage, and overall return on investment.
- Order Execution Logic: Simulates market conditions to determine if and when orders would be filled based on price movements.

**Data Visualization:**
- Price and Orders Plot: Generates a plot showing the asset's price over time along with the placement and execution of buy and sell orders.
- Spread Percentages Plot: Visualizes the spread percentages, indicating the potential profit margins at different times.
- Wallet Worth Plot: Displays the total value of the wallet over time, reflecting gains or losses from the trading activity.
- Image Saving: Saves all generated plots as image files (PNG format) for easy viewing and analysis.

**Logging and Reporting:**
- Detailed Logs: Provides console output of each iteration, including price updates, volatility predictions, order placements, and wallet status.
- Summary Statistics: Outputs final performance metrics after the simulation completes, allowing for quick assessment of strategy success.

## Configuration and Customization
- Configurable Parameters: Key settings such as the trading pair, GARCH model coefficients, order expiration, trade sizes, and wallet balances can be adjusted in `config.go`.
- Scalability: The simulation can be scaled by changing the `MaxIterations` parameter to run over longer periods or more data points.
- Strategy Adjustments: Traders can modify the trading logic, such as the z-score coefficient or order placement rules, to test different strategies.

## Setting Up and Running the Project
1. Install Go: Ensure Go is installed on your system (version 1.16 or later recommended).
2. Clone Repository: Download the project files into a directory on your machine.
3. Install Dependencies: Use `go mod tidy` to automatically download and install required packages.
4. Build Application: Run `go build` to compile the application.
5. Execute Simulation: Run the compiled binary to start the simulation. The application will connect to Binance, retrieve data, and perform the trading simulation.
6. View Results: After completion, open the generated image files (`prices.png`, `spread_percents.png`, `wallet_worth.png`) to analyze the visualizations.

## Note

This project is developed for educational purposes to demonstrate the implementation of a trading simulation using real-time market data and volatility modeling. Please be aware of the following theoretical issue in the current implementation:

- Data Input for GARCH Model: The GARCH model in this project uses ticker price updates as input data. However, GARCH models are theoretically designed to work with time series data that are equally spaced in time—such as closing prices from candlestick (OHLCV) data of a specific timeframe (e.g., 1-minute, 5-minute intervals).

- Implications: Using irregularly timed ticker data can lead to inaccurate volatility estimates because the time intervals between data points are inconsistent. This affects the reliability of the model's predictions and, consequently, the performance of the simulated trading strategy.

- Recommendation: For a theoretically sound application of the GARCH model, it is recommended to use evenly spaced time series data. Modifying the project to aggregate ticker data into candlestick data before applying the GARCH model would improve the theoretical correctness of the simulation.
