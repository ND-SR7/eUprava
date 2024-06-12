import { useNavigate } from "react-router-dom";
import Button from "../components/Shared/Button/Button";
import HeadingStyled from "../components/Shared/Heading/Heading.styled";
import { useEffect, useState } from "react";
import { pingPolice } from "../services/PingService";
import toast from "react-hot-toast";
import GetAllTrafficViolations from "../components/Police/TrafficViolation/GetAllTrafficViolation";
import Modal from "../components/Shared/Modal/Modal";

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

  const closeModal = () => {
    setIsModalVisible(false);
  };
  
  return (
    <>
      <HeadingStyled>Traffic Police</HeadingStyled>
      <br />
      <Button buttonType="button" label="Ping Service" onClick={() => ping()} />
      <br />
      <GetAllTrafficViolations setModalContent={setModalContent} setIsModalVisible={setIsModalVisible} />
      <br />
      <Button buttonType="button" label="Go Back" onClick={() => window.history.back()}/>
      <br />
      <Modal
        heading="Statistics"
        content={modalContent}
        isVisible={isModalVisible}
        onClose={closeModal}
      />
      <br />
    </>
  );
};

export default PolicePage;
