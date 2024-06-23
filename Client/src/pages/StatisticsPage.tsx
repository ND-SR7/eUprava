import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect, useState } from "react";
import { pingStatistics } from "../services/PingService";
import toast from "react-hot-toast";
import Modal from "../components/Shared/Modal/Modal";
import TrafficStatistics from "../components/Statistics/TrafficStatistics/TrafficStatistics";
import VehicleStatistics from "../components/Statistics/VehicleStatistics/VehicleStatistics";
import RegisteredVehicles from "../components/Statistics/RegisteredVehiclesByYear/RegisteredVehicles";

const StatisticsPage = () => {
  const navigate = useNavigate();
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

  const closeModal = () => {
    setIsModalVisible(false);
  };

  return (
    <>
      <HeadingStyled>Institute for Statistics</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <Button buttonType="button" label="Fetch Vehicle Statistics" onClick={() => {
        setModalContent(<VehicleStatistics />);
        setIsModalVisible(true);
      }} />
      <br />
      <Button buttonType="button" label="Fetch Traffic Statistics" onClick={() => {
        setModalContent(<TrafficStatistics />);
        setIsModalVisible(true);
      }} />
      <br />
      <Button buttonType="button" label="Search Registered Vehicles by Year" onClick={() => {
        setModalContent(<RegisteredVehicles />);
        setIsModalVisible(true);
      }} />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()} />
      <br />
      <Modal
        heading="Statistics"
        content={modalContent}
        isVisible={isModalVisible}
        onClose={closeModal}
      />
    </>
  );
};

export default StatisticsPage;
