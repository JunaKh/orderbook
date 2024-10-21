<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import Chart from 'chart.js/auto';

  let chart;
  let averagePrice: number = 0;
  let priceHistory: number[] = [];
  let ws: WebSocket;

  function createChart(ctx: CanvasRenderingContext2D) {
    return new Chart(ctx, {
      type: "line",
      data: {
        labels: [],
        datasets: [{
          label: "Average Price",
          data: [],
          borderColor: "rgba(75, 192, 192, 1)",
          backgroundColor: "rgba(75, 192, 192, 0.2)",
          borderWidth: 2,
          pointBackgroundColor: "#fff",
          pointBorderColor: "#000",
        }]
      },
      options: {
        layout: {
          padding: {
            bottom: 10,
          },
        },
        scales: {
          x: {
            ticks: { color: "#fff" },
          },
          y: {
            ticks: { color: "#fff" },
          }
        }
      }
    });
  }

  function connectWebSocket() {
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = (event: MessageEvent) => {
      try {
        averagePrice = JSON.parse(event.data);
        priceHistory.push(averagePrice);

        chart.data.labels.push(new Date().toLocaleTimeString());
        chart.data.datasets[0].data.push(averagePrice);
        chart.update();
      } catch (error) {
        console.error("Error processing WebSocket message:", error);
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
      console.warn("WebSocket connection closed");
    };
  }

  onMount(() => {
    const canvasElement = document.getElementById("priceChart") as HTMLCanvasElement;
    if (!canvasElement) {
      console.error("Canvas element not found!");
      return;
    }

    const ctx = canvasElement.getContext("2d");
    if (!ctx) {
      console.error("Failed to get 2D context from canvas");
      return;
    }

    chart = createChart(ctx);
    connectWebSocket();
  });

  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
</script>

<main>
  <h1>Average Price: {averagePrice.toFixed(2)}</h1>
  <div class="chart-container">
    <canvas id="priceChart"></canvas>
  </div>
</main>

<style>
  main {
    font-family: 'Roboto', sans-serif;
    background: linear-gradient(45deg, #2b7276, #122466);
    color: #fff;
    text-align: center;
    padding: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    box-sizing: border-box;
    margin: 0;
  }

  h1 {
    font-size: 3rem;
    font-weight: 700;
    margin-bottom: 20px;
    color: #fff;
  }

  .chart-container {
    width: 100%;
    max-width: 1200px;
    height: 400px;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0;
    margin: 0;
  }

  canvas {
    width: 100%;
    height: 100%;
    background: none;
    border-radius: 0;
    box-shadow: none;
  }

  @media (max-width: 768px) {
    .chart-container {
      width: 100%;
    }

    canvas {
      width: 100%;
      max-width: 100%;
    }
  }
</style>
