import { useState } from "react";
import { getVehicleStatisticsByYear } from "../../../services/StatisticsService";
import { VehicleStatisticsTable, NoVehicleStatisticsMessage, StyledButton } from "./VehicleStatistics.styled";

type VehicleStatisticsProps = {
  setModalContent: (content: React.ReactNode) => void;
  setIsModalVisible: (visible: boolean) => void;
};

const VehicleStatistics: React.FC<VehicleStatisticsProps> = ({ setModalContent, setIsModalVisible }) => {
  const [vehicleStatistics, setVehicleStatistics] = useState<{ [year: number]: number }>({});

  const fetchVehicleStatistics = () => {
    getVehicleStatisticsByYear()
      .then((data) => {
        setVehicleStatistics(data);
        setModalContent(
          <>
            <h2>Vehicle Statistics by Year</h2>
            {Object.keys(data).length > 0 ? (
              <VehicleStatisticsTable>
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
              </VehicleStatisticsTable>
            ) : (
              <NoVehicleStatisticsMessage>No statistics available</NoVehicleStatisticsMessage>
            )}
          </>
        );
        setIsModalVisible(true);
      })
      .catch((error) => {
        console.error(error);
      });
  };

  return <StyledButton onClick={fetchVehicleStatistics}>Fetch Vehicle Statistics</StyledButton>;
};

export default VehicleStatistics;
