class UseOcclusion(string filename = "", coordsys = "")
{
  public void surface(output color Ci, Oi)
  {
    normal Nn = normalize(N);
    float occ = 0;

    texture3d(filename, P, Nn, "coordsystem", coordsys,
              "_occlusion", occ);

    Ci = (1 - occ) * Cs * Os;
    Oi = Os;
  }
}
