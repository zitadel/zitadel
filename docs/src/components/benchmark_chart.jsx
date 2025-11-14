import React from "react";
import Chart from "react-google-charts";

export function BenchmarkChart({ testResults = [], height = '500px' } = {}) {
    if (!Array.isArray(testResults)) {
        console.error("BenchmarkChart: testResults is not an array. Received:", testResults);
        return <p>Error: Benchmark data is not available or in the wrong format.</p>;
    }

    if (testResults.length === 0) {
        return <p>No benchmark data to display.</p>;
    }

    const dataPerMetric = new Map();
    let maxVValue = 0;

    testResults.forEach((result) => {
        if (!result || typeof result.metric_name === 'undefined') {
            console.warn("BenchmarkChart: Skipping invalid result item:", result);
            return;
        }
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
        if (result.p99 !== undefined && result.p99 > maxVValue) {
            maxVValue = result.p99;
        }
        dataPerMetric.get(result.metric_name).push([
            result.timestamp ? new Date(result.timestamp) : null,
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
            maxValue: maxVValue > 0 ? maxVValue : undefined,
        },
        title: ''
    };
    const charts = [];
    
    dataPerMetric.forEach((data, metric) => {
        const opt = { ...options };
        opt.title = metric;
        charts.push(
            <Chart
                key={metric}
                chartType="LineChart"
                width="100%"
                height={height}
                options={opt}
                data={data}
                legendToggle
            />
        );
    });

    if (charts.length === 0) {
        return <p>No chart data could be generated.</p>;
    }

    return <>{charts}</>;
}