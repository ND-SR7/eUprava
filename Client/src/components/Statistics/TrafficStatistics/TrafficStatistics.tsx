import { useState } from "react";
import { getTrafficStatistics } from "../../../services/StatisticsService";
import { TrafficStatisticsTable, NoTrafficStatisticsMessage, StyledButton } from "./TrafficStatistics.styled";
import TrafficStatistic from "../../../models/Statistics/TrafficStatistic";
import React from "react";

type TrafficStatisticsProps = {
  setModalContent: (content: React.ReactNode) => void;
  setIsModalVisible: (visible: boolean) => void;
};

const TrafficStatistics: React.FC<TrafficStatisticsProps> = ({ setModalContent, setIsModalVisible }) => {
  const [trafficStatistics, setTrafficStatistics] = useState<TrafficStatistic[]>([]);

  const fetchTrafficStatistics = () => {
    getTrafficStatistics()
      .then((data) => {
        setTrafficStatistics(data);
        setModalContent(
          <>
            <h2>Traffic Statistics</h2>
            {data.length > 0 ? (
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
                  {data.map((stat: TrafficStatistic) => (
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
        setIsModalVisible(true);
      })
      .catch((error) => {
        console.error(error);
      });
  };

  return <StyledButton onClick={fetchTrafficStatistics}>Fetch Traffic Statistics</StyledButton>;
};

export default TrafficStatistics;
