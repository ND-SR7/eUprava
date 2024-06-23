import { useState, useEffect } from "react";
import { getTrafficStatistics } from "../../../services/StatisticsService";
import { TrafficStatisticsTable, NoTrafficStatisticsMessage, DownloadButton } from "./TrafficStatistics.styled";
import TrafficStatistic from "../../../models/Statistics/TrafficStatistic";
import React from "react";
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';

const TrafficStatistics: React.FC = () => {
  const [trafficStatistics, setTrafficStatistics] = useState<TrafficStatistic[]>([]);

  useEffect(() => {
    getTrafficStatistics()
      .then((data) => {
        setTrafficStatistics(data);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  const generatePDF = () => {
    const input = document.getElementById("traffic-statistics-table");
    if (input) {
      html2canvas(input).then((canvas) => {
        const imgData = canvas.toDataURL('image/png');
        const pdf = new jsPDF();
        const imgWidth = 190;
        const pageHeight = 295;
        const imgHeight = (canvas.height * imgWidth) / canvas.width;
        let heightLeft = imgHeight;

        let position = 10;

        pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
        heightLeft -= pageHeight;

        while (heightLeft >= 0) {
          position = heightLeft - imgHeight;
          pdf.addPage();
          pdf.addImage(imgData, 'PNG', 10, position, imgWidth, imgHeight);
          heightLeft -= pageHeight;
        }
        
        pdf.save("Traffic_Statistics_Report.pdf");
      });
    }
  };

  return (
    <>
    <div id="traffic-statistics-table">
      <h2>Traffic Statistics</h2>
      {trafficStatistics.length > 0 ? (
        <TrafficStatisticsTable>
          <thead>
            <tr>
              <th>Date</th>
              <th>Region</th>
              <th>Violation Type</th>
              <th>Vehicle Brand</th>
              <th>Vehicle Model</th>
              <th>Vehicle Year</th>
              <th>Registration</th>
              <th>Plates</th>
              <th>Owner</th>
            </tr>
          </thead>
          <tbody>
            {trafficStatistics.map((stat: TrafficStatistic) => (
              <tr key={stat.id}>
                <td>{new Date(stat.date).toLocaleDateString()}</td>
                <td>{stat.region}</td>
                <td>{stat.violationType}</td>
                {stat.vehicles.map((vehicle) => (
                  <React.Fragment key={vehicle.id}>
                    <td>{vehicle.brand}</td>
                    <td>{vehicle.model}</td>
                    <td>{vehicle.year}</td>
                    <td>{vehicle.registration}</td>
                    <td>{vehicle.plates}</td>
                    <td>{vehicle.owner}</td>
                  </React.Fragment>
                ))}
              </tr>
            ))}
          </tbody>
        </TrafficStatisticsTable>
      ) : (
        <NoTrafficStatisticsMessage>No traffic statistics available</NoTrafficStatisticsMessage>
      )}
      </div>
      <DownloadButton onClick={generatePDF}>Download PDF</DownloadButton>
    </>
  );
};

export default TrafficStatistics;
