import React from "react";
import Chart from "react-google-charts";

export function BenchmarkChart(testResults=[], height='500px') {

    const options = {
        legend: { position: 'bottom' },
        focusTarget: 'category',
        hAxis: {
            title: 'timestamp',
        },
        vAxis: {
            title: 'latency (ms)',
        },
    };

    const data = [
        [
            {type:"datetime", label: "timestamp"},
            {type:"number", label: "p50"},
            {type:"number", label: "p95"},
            {type:"number", label: "p99"},
        ],
    ]

    JSON.parse(testResults.testResults).forEach((result) => {
        data.push([
            new Date(result.timestamp),
            result.p50,
            result.p95,
            result.p99,
        ])
    });

    return (
        <Chart
            chartType="LineChart"
            width="100%"
            height="500px"
            options={options}
            data={data}
            legendToggle
        />
  );
}