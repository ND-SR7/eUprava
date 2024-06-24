import { useNavigate } from "react-router-dom";
import Button from "../../components/Shared/Button/Button";
import HeadingStyled from "../../components/Shared/Heading/Heading.styled";
import { useEffect, useState } from "react";
import { pingPolice } from "../../services/PingService";
import toast from "react-hot-toast";
import Modal from "../../components/Shared/Modal/Modal";
import GetAllTrafficViolations from "../../components/Police/TrafficViolation/GetAllTrafficViolation";
import CheckAlcoholLevelForm from "../../components/Police/TrafficViolation/CheckAlcoholLevelForm";
import CheckVehicleTireForm from "../../components/Police/TrafficViolation/CheckVehicleTireForm";
import CheckDriverBanForm from "../../components/Police/TrafficViolation/CheckDriverBanForm";
import CheckDriverPermitValidityForm from "../../components/Police/TrafficViolation/CheckDriverPermitValidityForm";

const PolicePage = () => {
  const navigate = useNavigate();
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalContent, setModalContent] = useState<any>(null);

  useEffect(() => {
    const token = localStorage.getItem("token") || "";
    if (token === "") navigate("/");
  });

  const ping = () => {
    pingPolice().then(() => {
      toast.success("Good connection with traffic police service");
    }).catch((error) => {
      toast.error("No connection to traffic police service");
      console.error(error);
    });
  };

  const openTrafficViolationsModal = () => {
    setModalContent(<GetAllTrafficViolations />);
    setIsModalVisible(true);
  };

  const openCheckAlcoholLevelModal = () => {
    setModalContent(<CheckAlcoholLevelForm closeModal={closeModal} />);
    setIsModalVisible(true);
  };

  const openCheckVehicleTireModal = () => {
    setModalContent(<CheckVehicleTireForm closeModal={closeModal} />);
    setIsModalVisible(true);
  };

  const openCheckDriverBanModal = () => {
    setModalContent(<CheckDriverBanForm closeModal={closeModal} />);
    setIsModalVisible(true);
  };

  const openCheckDriverPermitValidityModal = () => {
    setModalContent(<CheckDriverPermitValidityForm closeModal={closeModal} />);
    setIsModalVisible(true);
  };

  const closeModal = () => {
    setIsModalVisible(false);
  };

  return (
    <>
      <HeadingStyled>Traffic Police</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <Button buttonType="button" label="Get All Traffic Violations" onClick={openTrafficViolationsModal} />
      <br />
      <Button buttonType="button" label="Check Alcohol Level" onClick={openCheckAlcoholLevelModal} />
      <br />
      <Button buttonType="button" label="Check Vehicle Tire" onClick={openCheckVehicleTireModal} />
      <br />
      <Button buttonType="button" label="Check Driver Ban" onClick={openCheckDriverBanModal} />
      <br />
      <Button buttonType="button" label="Check Driver Permit Validity" onClick={openCheckDriverPermitValidityModal} />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()} />
      <br />
      <Modal
        heading="Traffic Violations"
        content={modalContent}
        isVisible={isModalVisible}
        onClose={closeModal}
      />
    </>
  );
};

export default PolicePage;
