import { useState, useEffect } from "react";
import { getTrafficStatistics } from "../../../services/StatisticsService";
import { TrafficStatisticsTable, NoTrafficStatisticsMessage } from "./TrafficStatistics.styled";
import TrafficStatistic from "../../../models/Statistics/TrafficStatistic";
import React from "react";

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

  return (
    <>
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
    </>
  );
};

export default TrafficStatistics;
