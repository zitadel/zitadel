import React from "react";
import Chart from "react-google-charts";

export function BenchmarkChart(testResults=[], height='500px') {
    const dataPerMetric = new Map();
    let maxVValue = 0;

    JSON.parse(testResults.testResults).forEach((result) => {
        if (!dataPerMetric.has(result.metric_name)) {
            dataPerMetric.set(result.metric_name, [
                [
                    {type:"datetime", label: "timestamp"},
                    {type:"number", label: "p50"},
                    {type:"number", label: "p95"},
                    {type:"number", label: "p99"},
                ],
            ]);
        }
        if (result.p99 > maxVValue) {
            maxVValue = result.p99;
        }
        dataPerMetric.get(result.metric_name).push([
            new Date(result.timestamp),
            result.p50,
            result.p95,
            result.p99,
        ]);
    });

    const options = {
        legend: { position: 'bottom' },
        focusTarget: 'category',
        hAxis: {
            title: 'timestamp',
        },
        vAxis: {
            title: 'latency (ms)',
            maxValue: maxVValue,
        },
        title: ''
    };
    const charts = [];
    
    dataPerMetric.forEach((data, metric) => {
        const opt = Object.create(options);
        opt.title = metric;
        charts.push(
            <Chart
                chartType="LineChart"
                width="100%"
                height={height}
                options={opt}
                data={data}
                legendToggle
            />
        );
    });


    return (charts);
}