import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect, useState } from "react";
import { pingStatistics } from "../services/PingService";
import { getVehicleStatisticsByYear, getTrafficStatistics } from "../services/StatisticsService";
import toast from "react-hot-toast";
import { StatisticsContainer, StatisticsTable, NoStatisticsMessage } from '../components/Statistics/StatisticsPage.styled';
import React from "react";
import Modal from "../components/Shared/Modal/Modal";
import { TrafficStatistics } from "../components/models/Statistic/TrafficStatistic";

type VehicleStatistics = {
  [year: number]: number;
};

const StatisticsPage = () => {
  const navigate = useNavigate();
  const [vehicleStatistics, setVehicleStatistics] = useState<VehicleStatistics>({});
  const [trafficStatistics, setTrafficStatistics] = useState<TrafficStatistics[]>([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalContent, setModalContent] = useState<any>(null);

  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");
  }, [navigate]);

  const ping = () => {
    pingStatistics().then(() => {
      toast.success("Good connection with Institute for Statistics service");
    }).catch((error) => {
      toast.error("No connection to Institute for Statistics service");
      console.error(error);
    });
  };

  const fetchVehicleStatistics = () => {
    getVehicleStatisticsByYear()
      .then((data) => {
        setVehicleStatistics(data);
        setModalContent(
          <>
            <h2>Vehicle Statistics by Year</h2>
            {Object.keys(data).length > 0 ? (
              <StatisticsTable>
                <thead>
                  <tr>
                    <th>Year</th>
                    <th>Count</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.entries(data).map(([year, count]) => (
                    <tr key={year}>
                      <td>{year}</td>
                      <td>{count}</td>
                    </tr>
                  ))}
                </tbody>
              </StatisticsTable>
            ) : (
              <NoStatisticsMessage>No statistics available</NoStatisticsMessage>
            )}
          </>
        );
        setIsModalVisible(true);
        toast.success("Vehicle statistics retrieved successfully");
      })
      .catch((error) => {
        toast.error("Failed to retrieve vehicle statistics");
        console.error(error);
      });
  };

  const fetchTrafficStatistics = () => {
    getTrafficStatistics()
      .then((data: TrafficStatistics[]) => {
        setTrafficStatistics(data);
        setModalContent(
          <>
            <h2>Traffic Statistics</h2>
            {data.length > 0 ? (
              <StatisticsTable>
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
                  {data.map((stat: TrafficStatistics) => (
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
              </StatisticsTable>
            ) : (
              <NoStatisticsMessage>No traffic statistics available</NoStatisticsMessage>
            )}
          </>
        );
        setIsModalVisible(true);
        toast.success("Traffic statistics retrieved successfully");
      })
      .catch((error) => {
        toast.error("Failed to retrieve traffic statistics");
        console.error(error);
      });
  };

  const closeModal = () => {
    setIsModalVisible(false);
  };

  return (
    <StatisticsContainer>
      <HeadingStyled>Institute for Statistics</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <Button buttonType="button" label="Fetch Vehicle Statistics" onClick={() => fetchVehicleStatistics()} />
      <br />
      <Button buttonType="button" label="Fetch Traffic Statistics" onClick={() => fetchTrafficStatistics()} />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()} />
      <br />
      <Modal
        heading="Statistics"
        content={modalContent}
        isVisible={isModalVisible}
        onClose={closeModal}
      />
    </StatisticsContainer>
  );
};

export default StatisticsPage;
